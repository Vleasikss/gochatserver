package main

import (
	"time"

	"github.com/Vleasikss/gochatserver/controllers"
	"github.com/Vleasikss/gochatserver/models"
	"github.com/Vleasikss/gochatserver/mongo"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"gopkg.in/olahol/melody.v1"
)

func main() {
	r := gin.Default()
	mongo := mongo.NewMongoClient[models.Message]()
	melody := melody.New()

	controller := controllers.NewMessageController(mongo, melody)

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5002"},
		AllowMethods:     []string{"PUT", "PATCH", "GET"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return true
		},
		MaxAge: 12 * time.Hour,
	}))

	r.Use(static.Serve("/", static.LocalFile("./public", true)))
	r.GET("/ws", controller.HandleSocketMessage)
	r.GET("/history", controller.FindMessageHistory)

	r.Run(":5002")
}
