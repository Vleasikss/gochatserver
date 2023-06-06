package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"

	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
)

type ClientManager struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

type Client struct {
	id     string
	socket *websocket.Conn
	send   chan []byte
}

type Message struct {
	Sender    string `json:"sender,omitempty"`
	Recipient string `json:"recipient,omitempty"`
	Content   string `json:"content,omitempty"`
	ServerIP  string `json:"serverIp,omitempty"`
	SenderIP  string `json:"senderIp,omitempty"`
}

var manager = ClientManager{
	broadcast:  make(chan []byte),
	register:   make(chan *Client),
	unregister: make(chan *Client),
	clients:    make(map[*Client]bool),
}

func (manager *ClientManager) start() {
	for {
		select {
		case conn := <-manager.register:
			manager.clients[conn] = true
			jsonMessage, _ := json.Marshal(&Message{
				Content:  "/A new socket has connect. ",
				ServerIP: LocalIp(),
				SenderIP: conn.socket.RemoteAddr().String(),
			})
			manager.send(jsonMessage, conn)

		}
	}
}

// Define the send method of client management
func (manager *ClientManager) send(message []byte, ignore *Client) {
	for conn := range manager.clients {
		if conn != ignore {
			//Send messages not to the shielded connection
			conn.send <- message
		}
	}
}

func (c *Client) read() {
	defer func() {
		manager.unregister <- c
		_ = c.socket.Close()
	}()

	for {
		_, message, err := c.socket.ReadMessage()
		fmt.Println("get message: " + string(message))
		//If there is an error message, cancel this connection and then close it
		if err != nil {
			manager.unregister <- c
			_ = c.socket.Close()
			break
		}
		//If there is no error message, put the information in Broadcast
		jsonMessage, _ := json.Marshal(&Message{
			Sender:   c.id,
			Content:  string(message),
			ServerIP: LocalIp(),
			SenderIP: c.socket.RemoteAddr().String(),
		})
		manager.broadcast <- jsonMessage
	}
}

func (c *Client) write() {
	defer func() {
		_ = c.socket.Close()
	}()

	for {
		select {
		//Read the message from send
		case message, ok := <-c.send:
			//If there is no message
			fmt.Println(string(message), ok)
			if !ok {
				_ = c.socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			//Write it if there is news and send it to the web side
			_ = c.socket.WriteMessage(websocket.TextMessage, message)
		}
	}
}

func main() {
	fmt.Println("Starting application...")
	//Open a goroutine execution start program
	go manager.start()
	//Register the default route to /ws, and use the wsHandler method
	http.HandleFunc("/ws", wsHandler)
	http.HandleFunc("/health", healthHandler)
	//Surveying the local 8011 port
	fmt.Println("chat server start.....")
	//Note that this must be 0.0.0.0 to deploy in the server to use
	_ = http.ListenAndServe("0.0.0.0:8448", nil)
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024 * 1024 * 1024,
	WriteBufferSize: 1024 * 1024 * 1024,
	//Solving cross-domain problems
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func wsHandler(res http.ResponseWriter, req *http.Request) {
	//Upgrade the HTTP protocol to the websocket protocol
	conn, err := upgrader.Upgrade(res, req, nil)
	if err != nil {
		fmt.Println(err)
		http.NotFound(res, req)
		return
	}

	//Every connection will open a new client, client.id generates through UUID to ensure that each time it is different
	client := &Client{
		id:     uuid.Must(uuid.NewV4(), nil).String(),
		socket: conn,
		send:   make(chan []byte),
	}
	//Register a new client
	manager.register <- client

	//Start the message to collect the news from the web side
	go client.read()
	//Start the corporation to return the message to the web side
	go client.write()
}

func healthHandler(res http.ResponseWriter, _ *http.Request) {
	_, _ = res.Write([]byte("ok"))
}

func LocalIp() string {
	address, _ := net.InterfaceAddrs()
	var ip = "localhost"
	for _, address := range address {
		if ipAddress, ok := address.(*net.IPNet); ok && !ipAddress.IP.IsLoopback() {
			if ipAddress.IP.To4() != nil {
				ip = ipAddress.IP.String()
			}
		}
	}
	return ip
}
