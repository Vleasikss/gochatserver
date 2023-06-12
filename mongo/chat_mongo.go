package mongo

import (
	"context"
	"fmt"
	"github.com/Vleasikss/gochatserver/models"
	"go.mongodb.org/mongo-driver/bson"
)

const ChatDatabase = "test"
const ChatCollection = "chats"

func (cl *MongoClient) FindAllChats() ([]models.Chat, error) {
	cl.mongo.Database(MessageDatabase).CreateCollection(context.Background(), ChatCollection)
	collection := cl.mongo.Database(ChatDatabase).Collection(ChatCollection)
	var res []models.Chat
	filter := bson.D{}
	cursor, err := collection.Find(context.Background(), filter)
	if err = cursor.All(context.Background(), &res); err != nil {
		return nil, err
	}
	cursor.Close(context.Background())

	return res, nil
}

func (cl *MongoClient) FindAllUserChats(userId uint) ([]models.Chat, error) {
	collection := cl.mongo.Database(ChatDatabase).Collection(ChatCollection)
	filter := bson.M{"assignedTo": userId}
	c, err := collection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	var res []models.Chat
	if err = c.All(context.Background(), &res); err != nil {
		return nil, err
	}

	c.Close(context.Background())

	return res, nil
}

func (cl *MongoClient) InsertChat(data *models.Chat) error {
	collection := cl.mongo.Database(ChatDatabase).Collection(ChatCollection)
	result, err := collection.InsertOne(context.Background(), &data)
	if err != nil {
		return err
	}
	fmt.Printf("Inserted document with _id: %v\n", result.InsertedID)
	return nil
}
