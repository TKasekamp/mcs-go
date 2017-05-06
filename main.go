package main

import (
	"log"
	"net/http"
	"os"
	"github.com/gin-gonic/gin"
	"fmt"
	"github.com/satori/go.uuid"
	"github.com/gorilla/websocket"
)

var commands []Command

//https://scotch.io/bar-talk/build-a-realtime-chat-server-with-go-and-websockets
var clients = make(map[*websocket.Conn]bool) // connected clients
var broadcast = make(chan Command)           // broadcast channel
// CheckOrigin needed because CORS or something
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	}, }

type Command struct {
	Id      string `json:"id"`
	Command string `json:"command"`
	UserId  string `json:"userId"`
	Status  string `json:status`
	Result  string `json:result`
}

func makeDB() []Command {
	//var c1 = Command{"rand-id-0", "print this", "random-user-id-0"}
	var commands []Command
	//commands = append(commands, c1)
	return commands
}

func processCommand(c *Command) {
	c.Id = fmt.Sprintf("%s", uuid.NewV4())
	c.Status = "Accepted"
}

func posting(c *gin.Context) {
	var json Command
	if c.BindJSON(&json) == nil {
		processCommand(&json)
		//fmt.Print(uuid.NewV4())
		//json.Id = fmt.Sprintf("%s", uuid.NewV4())
		commands = append(commands, json)
		broadcast <- json
		c.JSON(http.StatusOK, gin.H{"status": json.Status, "id": json.Id})
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
		var command Command
		// Read in a new message as JSON and map it to a Message object
		err := ws.ReadJSON(&command)
		if err != nil {
			log.Printf("error: %v", err)
			delete(clients, ws)
			break
		}
		// Send the newly received message to the broadcast channel
		broadcast <- command
	}
}

func handleMessages() {
	for {
		// Grab the next message from the broadcast channel
		msg := <-broadcast
		// Send it out to every client that is currently connected
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("error: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
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
	router.OPTIONS("/commands", OptionsUser) // POST
	router.GET("/commands", func(c *gin.Context) {
		c.JSON(http.StatusOK, commands)
	})

	router.POST("/commands", posting)

	// Configure websocket route
	router.GET("/ws", func(c *gin.Context) {
		handleConnections(c.Writer, c.Request)
	})

	// Start listening for incoming chat messages
	go handleMessages()

	router.Run(":" + port)
}
