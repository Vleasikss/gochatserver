package controllers

import (
	"fmt"
	"github.com/Vleasikss/gochatserver/jwt"
	"github.com/Vleasikss/gochatserver/models"
	"github.com/Vleasikss/gochatserver/mongo"
	"github.com/gin-gonic/gin"
	"gopkg.in/olahol/melody.v1"
	"log"
	"net/http"
)

type ChatController struct {
	MongoClient *mongo.Client
	Melody      *melody.Melody
}

func NewChatController(mongo *mongo.Client, m *melody.Melody) *ChatController {
	return &ChatController{
		MongoClient: mongo,
		Melody:      m,
	}
}

type NewChatRequest struct {
	Name         string   `json:"name" binding:"required"`
	Participants []string `json:"participants" binding:"required"`
}

func (cc *ChatController) PostChat(c *gin.Context) {

	var input NewChatRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userId, err := jwt.ExtractTokenID(c)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}
	user, _ := models.GetUserByID(userId)
	// admin is also participant

	chat := models.NewChat(input.Name, input.Participants, user.Username)
	log.Printf("posting new for user: userId=%d, chat: %s\n", userId, chat)

	err = cc.MongoClient.InsertChat(chat)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"chatId":  chat.ChatId,
		"message": "success",
	})
}

type DeleteChatRequest struct {
	ChatId string `json:"chatId" binding:"required"`
}

func (cc *ChatController) DeleteChat(c *gin.Context) {
	var input DeleteChatRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := cc.MongoClient.DeleteChatById(input.ChatId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
}

func (cc *ChatController) FindAllUserChats(c *gin.Context) {
	userId, _ := jwt.ExtractTokenID(c)
	user, _ := models.GetUserByID(userId)

	results, err := cc.MongoClient.FindAllUserChats(user)
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	fmt.Println(results)

	c.JSON(http.StatusOK, gin.H{"data": results})
}
