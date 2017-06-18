package command

import (
	"time"
	"fmt"
	"github.com/satori/go.uuid"
	"math/rand"
)

var statuses = []string{"RESPONSE_RECEIVED", "FAILED"}

type Command struct {
	Id          string `json:"id"`
	Body        string `json:"body"`
	PrototypeId int `json:"prototypeId"`
	Priority    string `json:"priority"`
	UserId      string `json:"userId"`
	Status      string `json:"status"`
	Response    Response `json:"response"`
	SubmitTime  time.Time `json:"submitTime"`
}

type Response struct {
	Body         map[string]string `json:"body"`
	ResponseTime time.Time `json:"responseTime"`
}

func (c *Command) ProcessCommand() {
	c.Id = fmt.Sprintf("%s", uuid.NewV4())
	c.Status = "ACCEPTED"
	c.UserId = "server-value"
	c.SubmitTime = time.Now().UTC()
}

func (msg *Command) SimulateWork() {
	fmt.Println("Doing some work")
	//Simulate work
	if msg.Body != "a" {
		time.Sleep(time.Second * 3)
	}
	msg.Status = statuses[random(0, len(statuses))]
	if msg.Status == "FAILED" {
		msg.Response.Body = map[string]string{"error": "burning!"}
	} else {
		msg.Response.Body = map[string]string{"apple": "5", "lettuce": "7"}
	}
	msg.Response.ResponseTime = time.Now().UTC()

}

func random(min, max int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(max-min) + min
}
