package mongo

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Replace the placeholder with your Atlas connection string
const uri = "mongodb://mongodb:27017"

type Client struct {
	mongo *mongo.Client
}

func connect() (*mongo.Client, error) {
	// Use the SetServerAPIOptions() method to set the Stable API version to 1
	// serverAPI := options.ServerAPI(options.ServerAPIVersion1)

	fmt.Println("Connecting to database...")
	opts := options.Client().ApplyURI(uri)

	// Create a new client and connect to the server
	connection, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		fmt.Println("unable to connect to mongodb: " + err.Error())
		return nil, err
	}
	// Check the connection
	err = connection.Ping(context.TODO(), nil)
	if err != nil {
		fmt.Println("unable to connect to mongodb: " + err.Error())
	}
	return connection, nil
}

func NewMongoClient() *Client {
	cl, err := connect()
	if err != nil {
		fmt.Println("error during mongo connection: " + err.Error())
		return nil
	}
	return &Client{
		mongo: cl,
	}
}

func (cl Client) Disconnect() {
	if err := cl.mongo.Disconnect(context.TODO()); err != nil {
		fmt.Println("Unable to disconnect: " + err.Error())
	}
}
