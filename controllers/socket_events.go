package controllers

import "encoding/json"

type EventType string

const (
	MessageCreated EventType = "MESSAGE_CREATED"
	ChatCreated    EventType = "CHAT_CREATED"
	ChatDeleted    EventType = "CHAT_DELETED"
)

type RawSocketMessage struct {
	EventType EventType        `json:"eventType"`
	Message   *json.RawMessage `json:"message"`
}
type RawSocketEvent struct {
	EventType EventType `json:"eventType"`
}

type SocketMessage struct {
	EventType EventType   `json:"eventType"`
	Message   interface{} `json:"message"`
}

type MessageCreatedEvent struct {
	From     string `json:"from"`
	Payload  string `json:"payload"`
	ChatId   string `json:"chatId"`
	ChatName string `json:"chatName"`
}

type ChatCreatedEvent struct {
	ChatId       string   `json:"chatId"`
	Name         string   `json:"name"`
	AssignedTo   string   `json:"assignedTo"`
	Participants []string `json:"participants"`
}

func NewChatCreatedEvent(chatId string, name string, assignedTo string, participants []string) *SocketMessage {
	return &SocketMessage{
		ChatCreated,
		ChatCreatedEvent{
			ChatId:       chatId,
			Name:         name,
			AssignedTo:   assignedTo,
			Participants: participants,
		},
	}
}

type ChatDeletedEvent struct {
	ChatId string `json:"chatId"`
}
