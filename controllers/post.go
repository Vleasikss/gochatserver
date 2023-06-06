package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"samzhangjy/go-blog/mongo"
	
	"github.com/gin-gonic/gin"
	"gopkg.in/olahol/melody.v1"
)

type MessageController struct {
	MongoClient *mongo.MongoClient[Message]
	Melody      *melody.Melody
}

func NewMessageController(mongo *mongo.MongoClient[Message], m *melody.Melody) *MessageController {

	m.HandleMessage(func(s *melody.Session, msg []byte) {
		var input Message
		err := json.Unmarshal(msg, &input)
		if err != nil {
			fmt.Println("error during JSON parsing: " + err.Error())
		}
		fmt.Printf("Inserting message: from=%s, payload=%s", input.From, input.Payload)
		go mongo.Insert(&input)
		go m.Broadcast(msg)
	})

	return &MessageController{
		MongoClient: mongo,
		Melody:      m,
	}
}

func (mc *MessageController) FindMessageHistory(c *gin.Context) {
	fmt.Println("GET request to get the history. Started...")
	results := mc.MongoClient.FindAll()
	c.JSON(http.StatusOK, gin.H{"data": results})
}

func (mc *MessageController) HandleSocketMessage(c *gin.Context) {
	mc.Melody.HandleRequest(c.Writer, c.Request)
}
