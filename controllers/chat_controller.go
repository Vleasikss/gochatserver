package controllers

import (
	"fmt"
	"github.com/Vleasikss/gochatserver/jwt"
	"github.com/Vleasikss/gochatserver/models"
	"github.com/Vleasikss/gochatserver/mongo"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ChatController struct {
	MongoClient *mongo.MongoClient
}

func NewChatController(mongo *mongo.MongoClient) *ChatController {
	return &ChatController{
		MongoClient: mongo,
	}
}

type NewChatRequest struct {
	Name         string
	Participants []string
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

	chat := models.NewChat(input.Name, input.Participants, userId)

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

func (cc *ChatController) FindAllUserChats(c *gin.Context) {
	userId, _ := jwt.ExtractTokenID(c)

	results, err := cc.MongoClient.FindAllUserChats(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	fmt.Println(results)

	c.JSON(http.StatusOK, gin.H{"data": results})
}
