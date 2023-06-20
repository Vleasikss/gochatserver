package main

import (
	"fmt"
	"github.com/Vleasikss/gochatserver/controllers"
	"github.com/Vleasikss/gochatserver/jwt"
	"github.com/Vleasikss/gochatserver/models"
	"github.com/Vleasikss/gochatserver/mongo"
	"github.com/gin-contrib/cors"
	"os"
	"time"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"gopkg.in/olahol/melody.v1"
)

func main() {
	port := os.Getenv("port")
	r := gin.Default()
	models.ConnectDataBase()
	mongoClient := mongo.NewMongoClient()
	melodyClient := melody.New()

	mc := controllers.NewMessageController(mongoClient, melodyClient)
	cc := controllers.NewChatController(mongoClient, melodyClient)
	msc := controllers.NewMessageSocketController(mongoClient, melodyClient)

	r.Use(corsRules(port))
	r.Use(static.Serve("/", static.LocalFile("./public", true)))

	public := r.Group("/api")
	public.POST("/register", controllers.Authenticate)
	public.POST("/login", controllers.Login)
	public.GET("/ws", msc.HandleSocketMessage)

	protected := r.Group("/api")
	protected.Use(jwt.AuthMiddleware())
	protected.GET("/history/:chatId", mc.FindMessageHistory)
	protected.GET("/users", mc.FindAllUsers)
	protected.POST("/chat", cc.PostChat)
	protected.GET("/chat", cc.FindAllUserChats)
	protected.DELETE("/chat/:chatId", cc.DeleteChat)

	private := protected.Group("/admin")
	private.GET("/user", controllers.CurrentUser)

	err := r.Run(":" + port)
	if err != nil {
		fmt.Println("unexpected error during running application: " + err.Error())
	}
}

func corsRules(port string) gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost", "http://localhost:" + port, "http://localhost:3000"},
		AllowMethods:     []string{"PUT", "PATCH", "GET", "POST", "DELETE"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return true
		},
		MaxAge: 12 * time.Hour,
	})
}
