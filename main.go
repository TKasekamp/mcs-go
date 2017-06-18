package main

import (
	"log"
	"net/http"
	"os"
	"github.com/gin-gonic/gin"
	"fmt"
	"github.com/gorilla/websocket"
	"time"
	"mcs-go/command"
	"io/ioutil"
	"encoding/json"
)

var commands []command.Command
//https://scotch.io/bar-talk/build-a-realtime-chat-server-with-go-and-websockets
var clients = make(map[*websocket.Conn]bool) // connected clients
var workChannel = make(chan command.Command) // First channel
var broadcast = make(chan command.Command)   // broadcast channel
// CheckOrigin needed because CORS or something
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	}, }

func makeDB() []command.Command {
	//var c1 = Command{"rand-id-0", "print this", "random-user-id-0"}
	var commands []command.Command
	//commands = append(commands, c1)
	return commands
}

func posting(c *gin.Context) {
	var json command.Command
	if c.BindJSON(&json) == nil {
		(&json).ProcessCommand()
		//fmt.Print(uuid.NewV4())
		//json.Id = fmt.Sprintf("%s", uuid.NewV4())
		commands = append(commands, json)
		// Putting the command in some queue
		if json.Body != "a" {
			time.Sleep(time.Millisecond * 500)
		}
		c.JSON(http.StatusAccepted, json)
		workChannel <- json

	}
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	// Upgrade initial GET request to a websocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	// Make sure we close the connection when the function returns
	defer ws.Close()

	// Register our new client
	clients[ws] = true

	for {
		// I guess this has to be here???
		// There is actually no command sent here
		var c command.Command
		// Read in a new message as JSON and map it to a Message object
		err := ws.ReadJSON(&c)
		if err != nil {
			log.Printf("error: %v", err)
			delete(clients, ws)
			break
		}
		//// Send the newly received message to the broadcast channel
		//broadcast <- command
	}
}

func handleMessages() {
	for {
		// Grab the next message from the broadcast channel
		msg := <-broadcast

		// Send it out to every client that is currently connected
		for client := range clients {
			fmt.Println("Sending message to client")
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("error: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}

func handleWork() {
	for {
		msg := <-workChannel
		(&msg).SimulateWork()
		for i := range commands {
			if commands[i].Id == msg.Id {
				commands[i] = msg
				break
			}
		}
		broadcast <- msg
	}
}

func handleLoadJson(c *gin.Context) {
	//c.Writer.Header().Set("Content-Type", "application/json")
	//c.Next()

	//c.File("static/prototypes.json")
	//http.ServeFile(c.Writer, c.Request, "static/prototypes.json" )
	//c.Next()
	b, err := ioutil.ReadFile("static/prototypes.json")
	if err != nil {
		fmt.Println(err.Error())
	}
	var protos []command.Prototype
	json.Unmarshal(b, &protos)
	c.JSON(http.StatusOK, protos)
}

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "DELETE,POST, PUT")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		c.Next()
	}
}

func OptionsUser(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Methods", "DELETE,POST, PUT")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	c.Next()
}

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		// Pure development thing
		port = "8080"
		log.Print("$PORT must be set")
		//log.info("$PORT must be set")
	}

	commands = makeDB()

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(Cors())
	router.OPTIONS("/api/commands", OptionsUser) // POST
	router.GET("/api/commands", func(c *gin.Context) {
		c.JSON(http.StatusOK, commands)
	})

	router.POST("/api/commands", posting)
	router.GET("/api/prototypes", handleLoadJson)
	// Configure websocket route
	router.GET("/ws", func(c *gin.Context) {
		handleConnections(c.Writer, c.Request)
	})

	// Thread for working
	go handleWork()

	// Start listening for incoming chat messages
	go handleMessages()

	router.Run(":" + port)
}
