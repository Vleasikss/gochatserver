package controllers

import (
	"gopkg.in/olahol/melody.v1"
	"net/http"

	"github.com/Vleasikss/gochatserver/models"
	"github.com/Vleasikss/gochatserver/mongo"

	"github.com/gin-gonic/gin"
)

type MessageController struct {
	MongoClient *mongo.Client
	Melody      *melody.Melody
}

func NewMessageController(mongo *mongo.Client, melody *melody.Melody) *MessageController {
	return &MessageController{
		MongoClient: mongo,
		Melody:      melody,
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
