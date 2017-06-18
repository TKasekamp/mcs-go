package command

import "time"

type Command struct {
	Id             string `json:"id"`
	Body           string `json:"body"`
	PrototypeId	   int `json:"prototypeId"`
	Priority       string `json:"priority"`
	UserId         string `json:"userId"`
	Status         string `json:"status"`
	Response       Response `json:"response"`
	SubmitTime     time.Time `json:"submitTime"`
}

type Response struct {
	Body           map[string]string `json:"body"`
	ResponseTime   time.Time `json:"responseTime"`
}
