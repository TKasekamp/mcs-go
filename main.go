package main

import (
	"log"
	"net/http"
	"os"
	"github.com/gin-gonic/gin"
	"fmt"
	"github.com/satori/go.uuid"
)

var commands []Command

type Command struct {
	Id string `json:"id"`
	Command string `json:"command"`
	UserId string `json:"userId"`
}

func makeDB() []Command {
	var c1 = Command{"rand-id-0", "print this", "sda"}
	var commands []Command
	commands = append(commands, c1)
	return commands
}

func posting(c *gin.Context) {
	var json Command
	if c.BindJSON(&json) == nil {
		fmt.Print(uuid.NewV4())
		json.Id = fmt.Sprintf("%s", uuid.NewV4())
		commands = append(commands, json)
		c.JSON(http.StatusOK, gin.H{"status": "command submitted"})
	}
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

	router.GET("/commands", func(c *gin.Context) {
		c.JSON(http.StatusOK, commands)
	})

	router.POST("/commands", posting)


	router.Run(":" + port)
}
