package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func ws(c *gin.Context) {

	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Print("upgrade error: ", err)
		return
	}

	defer ws.Close()
	for {
		mt, message, err := ws.ReadMessage()
		if err != nil {
			log.Println("read error:", err)
			break
		}
		log.Printf("recv: %s\n", message)
		err = ws.WriteMessage(mt, message)
		if err != nil {
			log.Println("write error: ", err)
			break
		}
	}
}

func main() {
	fmt.Println("WebSocket Server!")
	bindAddress := "localhost:8448"
	r := gin.Default()
	r.GET("/ws", ws)
	r.Run(bindAddress)
}
