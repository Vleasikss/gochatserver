package controllers

import (
	"encoding/json"
	"github.com/Vleasikss/gochatserver/models"
	"github.com/Vleasikss/gochatserver/mongo"
	"github.com/gin-gonic/gin"
	"gopkg.in/olahol/melody.v1"
	"log"
)

type MessageSocketController struct {
	MongoClient *mongo.Client
	Melody      *melody.Melody
}

func NewMessageSocketController(mongo *mongo.Client, m *melody.Melody) *MessageSocketController {
	m.HandleMessage(onMessageReceived(mongo, m))

	return &MessageSocketController{
		MongoClient: mongo,
		Melody:      m,
	}
}

func handleMessageCreated(rawJson []byte, mongo *mongo.Client, m *melody.Melody) {
	type MessageCreated struct {
		EventType EventType           `json:"eventType"`
		Message   MessageCreatedEvent `json:"message"`
	}
	var msg MessageCreated
	if err := json.Unmarshal(rawJson, &msg); err != nil {
		log.Fatal(err)
		return
	}
	input := msg.Message
	event := MessageCreatedEvent{
		ChatId:   input.ChatId,
		From:     input.From,
		Payload:  input.Payload,
		ChatName: input.ChatName,
	}
	response := SocketMessage{
		EventType: "MESSAGE_CREATED",
		Message:   event,
	}
	message := models.Message{
		From:    input.From,
		Payload: input.Payload,
		ChatId:  input.ChatId,
	}
	log.Printf("Inserting message: from=%s, payload=%s", input.From, input.Payload)
	go mongo.InsertMessage(&message)
	responseJson, _ := json.Marshal(response)

	go m.Broadcast(responseJson)
}
func handleChatCreated(rawJson []byte, mongo *mongo.Client, m *melody.Melody) {
	type ChatCreated struct {
		EventType EventType        `json:"eventType"`
		Message   ChatCreatedEvent `json:"message"`
	}
	var message ChatCreated
	if err := json.Unmarshal(rawJson, &message); err != nil {
		log.Fatal(err)
	}
	log.Printf("ChatCreated event: %s\n", message)
	response := SocketMessage{
		EventType: "CHAT_CREATED",
		Message:   message.Message,
	}
	j, _ := json.Marshal(response)
	go m.Broadcast(j)
}

func handleChatDeleted(rawJson []byte, mongo *mongo.Client, m *melody.Melody, s *melody.Session) {
	type ChatDeleted struct {
		EventType EventType        `json:"eventType"`
		Message   ChatDeletedEvent `json:"message"`
	}
	var event ChatDeleted
	if err := json.Unmarshal(rawJson, &event); err != nil {
		log.Fatal(err)
	}
	log.Printf("sending event of recent deleted chat for event: %s\n", event)
	response := SocketMessage{
		EventType: "CHAT_DELETED",
		Message:   event.Message,
	}
	j, _ := json.Marshal(response)
	go m.Broadcast(j)
}

func onMessageReceived(mongo *mongo.Client, m *melody.Melody) func(s *melody.Session, msg []byte) {
	return func(s *melody.Session, msg []byte) {
		var event RawSocketEvent
		log.Printf("receive message: %s\n", string(msg))
		if err := json.Unmarshal(msg, &event); err != nil {
			log.Fatal(err)
			return
		}
		switch event.EventType {
		case MessageCreated:
			handleMessageCreated(msg, mongo, m)
		case ChatCreated:
			handleChatCreated(msg, mongo, m)
		case ChatDeleted:
			handleChatDeleted(msg, mongo, m, s)
		}
	}
}

func (mc *MessageSocketController) HandleSocketMessage(c *gin.Context) {
	mc.Melody.HandleRequest(c.Writer, c.Request)
}
