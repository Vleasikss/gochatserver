package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Vleasikss/gochatserver/models"
	"github.com/Vleasikss/gochatserver/mongo"

	"github.com/gin-gonic/gin"
	"gopkg.in/olahol/melody.v1"
)

type MessageController struct {
	MongoClient *mongo.MongoClient
	Melody      *melody.Melody
}

func NewMessageController(mongo *mongo.MongoClient, m *melody.Melody) *MessageController {

	m.HandleMessage(func(s *melody.Session, msg []byte) {
		var input models.Message
		err := json.Unmarshal(msg, &input)
		if err != nil {
			fmt.Println("error during JSON parsing: " + err.Error())
		}
		fmt.Printf("Inserting message: from=%s, payload=%s", input.From, input.Payload)
		go mongo.InsertMessage(&input)
		go m.Broadcast(msg)
	})

	return &MessageController{
		MongoClient: mongo,
		Melody:      m,
	}
}

func (mc *MessageController) FindMessageHistory(c *gin.Context) {
	chatId := c.Param("chatId")

	results, err := mc.MongoClient.FindAllMessagesByChatId(chatId)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"data": results})
}

func (mc *MessageController) FindAllUsers(c *gin.Context) {
	results, err := models.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": results})
}

func (mc *MessageController) HandleSocketMessage(c *gin.Context) {
	mc.Melody.HandleRequest(c.Writer, c.Request)
}
