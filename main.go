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
	mongoClient := mongo.NewMongoClient[models.Message]()
	melody := melody.New()

	controller := controllers.NewMessageController(mongoClient, melody)

	r.Use(corsRules(port))
	r.Use(static.Serve("/", static.LocalFile("./public", true)))

	public := r.Group("/api")
	public.POST("/register", controllers.Authenticate)
	public.POST("/login", controllers.Login)

	protected := r.Group("/api")
	protected.Use(jwt.AuthMiddleware())
	protected.GET("/ws", controller.HandleSocketMessage)
	protected.GET("/history", controller.FindMessageHistory)

	private := protected.Group("/admin")
	private.GET("/user", controllers.CurrentUser)

	err := r.Run(":" + port)
	if err != nil {
		fmt.Println("unexpected error during running application: " + err.Error())
	}
}

func corsRules(port string) gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost", "http://localhost:" + port},
		AllowMethods:     []string{"PUT", "PATCH", "GET"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return true
		},
		MaxAge: 12 * time.Hour,
	})
}
