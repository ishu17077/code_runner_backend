package database

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func dBInstance() *mongo.Client {
	mongoDb := "mongodb://localhost:27071"
	fmt.Print("Trying to connect to mongo db")

	fmt.Printf("Connecting to MongoDB:%s", mongoDb)

	client, err := mongo.Connect(options.Client().SetConnectTimeout(time.Second * 10).ApplyURI(mongoDb))

	if err != nil {
		fmt.Print("Connection to mongoDb failed")
		panic(err)
	}
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	fmt.Print("Successfully connected to mongoDB instance")
	return client
}

var Client *mongo.Client = dBInstance()

func OpenCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	var collection *mongo.Collection = client.Database("code_runner").Collection(collectionName)
	return collection
}
