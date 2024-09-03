package db

import (
	"context"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

func ConnectDB() {
	godotenv.Load()
	MONGO_URI := os.Getenv("MONGO_URI")
	configOptions := options.Client().ApplyURI(MONGO_URI)

	client, err := mongo.Connect(context.TODO(), configOptions)
	if err != nil {
		panic(err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		panic(err)
	}
	fmt.Println("Pinged successfully")
}

func DisconnectDB() {
	if client != nil {
		client.Disconnect(context.TODO())
	}
	fmt.Println("Closed DB")
}
