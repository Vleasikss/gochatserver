package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"samzhangjy/go-blog/mongo"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"gopkg.in/olahol/melody.v1"
)

type Message struct {
	From    string `json:"from"`
	Payload string `json:"payload"`
}

func main() {
	r := gin.Default()
	m := melody.New()
	cl := mongo.NewMongoClient[Message]()

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

	r.GET("/ws", func(c *gin.Context) {
		m.HandleRequest(c.Writer, c.Request)
	})
	r.GET("/history", func(c *gin.Context) {
		fmt.Println("GET request to get the history. Started...")
		results := cl.FindAll()
		c.JSON(http.StatusOK, gin.H{"data": results})
	})
	r.GET("/", func(c *gin.Context) {
		markdown, err := os.ReadFile("/app/go-sample-app/index.html")
		if err != nil {
			fmt.Println("error during reading /app/index.html: " + err.Error())
		}
		c.Data(http.StatusOK, "text/html; charset=utf-8", markdown)
	})

	m.HandleMessage(func(s *melody.Session, msg []byte) {
		var input Message
		err := json.Unmarshal(msg, &input)
		if err != nil {
			fmt.Println("error during JSON parsing: " + err.Error())
		}
		fmt.Printf("Inserting message: from=%s, payload=%s", input.From, input.Payload)
		go cl.Insert(&input)
		go m.Broadcast(msg)
	})

	r.Run(":5002")
}
