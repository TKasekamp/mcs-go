package main

import (
	"log"
	"net/http"
	"os"
	"github.com/gin-gonic/gin"
	"fmt"
	"github.com/satori/go.uuid"
	"github.com/gorilla/websocket"
	"math/rand"
	"time"
)

var commands []Command
var statuses = []string{"RESPONSE_RECEIVED", "FAILED"}
//https://scotch.io/bar-talk/build-a-realtime-chat-server-with-go-and-websockets
var clients = make(map[*websocket.Conn]bool) // connected clients
var workChannel = make(chan Command)         // First channel
var broadcast = make(chan Command)           // broadcast channel
// CheckOrigin needed because CORS or something
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	}, }

type Command struct {
	Id             string `json:"id"`
	Command        string `json:"commandString"`
	ObcsSchedule   string `json:"obcsSchedule"`
	McsSchedule    string `json:"mcsSchedule"`
	Priority       string `json:"priority"`
	UserId         string `json:"userId"`
	Status         string `json:"status"`
	ResponseString string `json:"responseString"`
	SubmitTime     time.Time `json:"submitTime"`
	ResponseTime   time.Time `json:"responseTime"`
}

func makeDB() []Command {
	//var c1 = Command{"rand-id-0", "print this", "random-user-id-0"}
	var commands []Command
	//commands = append(commands, c1)
	return commands
}

func processCommand(c *Command) {
	c.Id = fmt.Sprintf("%s", uuid.NewV4())
	c.Status = "ACCEPTED"
	c.SubmitTime = time.Now().UTC()
}

func posting(c *gin.Context) {
	var json Command
	if c.BindJSON(&json) == nil {
		processCommand(&json)
		//fmt.Print(uuid.NewV4())
		//json.Id = fmt.Sprintf("%s", uuid.NewV4())
		commands = append(commands, json)
		// Putting the command in some queue
		if json.Command != "a" {
			time.Sleep(time.Millisecond * 500)
		}
		c.JSON(http.StatusAccepted, json)
		workChannel <- json
		//c.JSON(http.StatusOK, gin.H{"status": json.Status, "id": json.Id})

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
		var command Command
		// Read in a new message as JSON and map it to a Message object
		err := ws.ReadJSON(&command)
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
		simulateWork(&msg)
		for i := range commands {
			if commands[i].Id == msg.Id {
				commands[i] = msg
				break
			}
		}
		broadcast <- msg
	}
}
func simulateWork(msg *Command) {
	fmt.Println("Doing some work")
	//Simulate work
	if msg.Command != "a" {
		time.Sleep(time.Second * 3)
	}
	msg.Status = statuses[random(0, len(statuses))]
	if msg.Status == "FAILED" {
		msg.ResponseString = "Something went very wrong"
	} else {
		msg.ResponseString = "Command queued up for next pass"
	}
	msg.ResponseTime = time.Now().UTC()

}

func random(min, max int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(max-min) + min
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

	// Thread for working
	go handleWork()

	// Start listening for incoming chat messages
	go handleMessages()

	router.Run(":" + port)
}
