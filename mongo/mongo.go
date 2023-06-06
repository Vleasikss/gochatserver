package mongo

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Replace the placeholder with your Atlas connection string
const uri = "mongodb://mongodb:27017"

type MongoClient[T any] struct {
	mongo *mongo.Client
}

func connect() (*mongo.Client, error) {
	// Use the SetServerAPIOptions() method to set the Stable API version to 1
	// serverAPI := options.ServerAPI(options.ServerAPIVersion1)

	fmt.Println("Connecting to database...")
	opts := options.Client().ApplyURI("mongodb://mongodb:27017")

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

func NewMongoClient[T any]() *MongoClient[T] {
	cl, err := connect()
	if err != nil {
		fmt.Println("error during mongo connection: " + err.Error())
		return nil
	}
	return &MongoClient[T]{
		mongo: cl,
	}
}

func (cl *MongoClient[T]) Insert(data *T) {
	collection := cl.mongo.Database("test").Collection("books")
	result, err := collection.InsertOne(context.Background(), &data)
	if err != nil {
		fmt.Println("error during insert: " + err.Error())
	}
	fmt.Printf("Inserted document with _id: %v\n", result.InsertedID)
}

func (cl *MongoClient[T]) FindAll() []T {
	cl.mongo.Database("test").CreateCollection(context.TODO(), "books")
	collection := cl.mongo.Database("test").Collection("books")
	fmt.Println("get the collection test/books: " + collection.Name())
	var res []T
	filter := bson.D{}
	cursor, err := collection.Find(context.TODO(), filter)
	if err = cursor.All(context.TODO(), &res); err != nil {
		fmt.Println("error during find all: " + err.Error())
	}
	fmt.Printf("get the data: %v\n", res)
	cursor.Close(context.TODO())

	return res
}

func (cl MongoClient[T]) Disconnect(data *T) {
	if err := cl.mongo.Disconnect(context.TODO()); err != nil {
		fmt.Println("Unable to disconnect: " + err.Error())
	}
}
