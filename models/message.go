package models

// Message mongo
type Message struct {
	From    string `json:"from" bson:"from"`
	Payload string `json:"payload" bson:"payload"`
	ChatId  string `json:"chatId" bson:"chatId"`
}
