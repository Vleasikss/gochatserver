package models


type Message struct {
	From    string `json:"from"`
	Payload string `json:"payload"`
}