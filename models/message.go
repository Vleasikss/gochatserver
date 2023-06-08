package models

// Message mongo
type Message struct {
	From    string `json:"from"`
	Payload string `json:"payload"`
}
