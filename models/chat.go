package models

import "github.com/google/uuid"

type Chat struct {
	ChatId       string   `json:"chatId" bson:"chatId"`
	Name         string   `json:"name" bson:"name"`
	Participants []string `json:"participants" bson:"participants"`
	AssignedTo   uint     `json:"assignedTo" bson:"assignedTo"`
}

func (c Chat) ExistingChat(chatId string, name string, participants []string, assignedTo uint) *Chat {
	return &Chat{
		chatId, name, participants, assignedTo,
	}
}

func NewChat(name string, participants []string, assignedTo uint) *Chat {
	return &Chat{
		ChatId:       generateChatId(),
		Name:         name,
		Participants: participants,
		AssignedTo:   assignedTo,
	}
}

func generateChatId() string {
	id := uuid.New()
	return id.String()
}
