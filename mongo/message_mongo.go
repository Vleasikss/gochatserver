package mongo

import (
	"context"
	"fmt"
	"github.com/Vleasikss/gochatserver/models"
	"go.mongodb.org/mongo-driver/bson"
)

const MessageDatabase = "test"
const MessageCollection = "messages"

func (cl *Client) InsertMessage(data *models.Message) {
	collection := cl.mongo.Database("test").Collection("books")
	result, err := collection.InsertOne(context.Background(), &data)
	if err != nil {
		fmt.Println("error during insert: " + err.Error())
	}
	fmt.Printf("Inserted document with _id: %v\n", result.InsertedID)
}

func (cl *Client) FindAllMessages() []models.Message {
	cl.mongo.Database(MessageDatabase).CreateCollection(context.Background(), MessageCollection)
	collection := cl.mongo.Database("test").Collection("books")
	fmt.Println("get the collection test/books: " + collection.Name())
	var res []models.Message
	filter := bson.D{}
	cursor, err := collection.Find(context.Background(), filter)
	if err = cursor.All(context.Background(), &res); err != nil {
		fmt.Println("error during find all: " + err.Error())
	}
	fmt.Printf("get the data: %v\n", res)
	cursor.Close(context.Background())

	return res
}

func (cl *Client) FindAllMessagesByChatId(chatId string) ([]models.Message, error) {
	collection := cl.mongo.Database("test").Collection("books")
	filter := bson.M{"chatId": chatId}

	c, err := collection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	var res []models.Message
	if err = c.All(context.Background(), &res); err != nil {
		return nil, err
	}
	c.Close(context.Background())

	return res, nil
}
