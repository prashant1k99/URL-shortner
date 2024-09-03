package db

import (
	"context"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoClient *mongo.Client

func ConnectDB() {
	godotenv.Load()
	MONGO_URI := os.Getenv("MONGO_URI")
	configOptions := options.Client().ApplyURI(MONGO_URI)

	var err error
	mongoClient, err = mongo.Connect(context.TODO(), configOptions)
	if err != nil {
		panic(err)
	}

	err = mongoClient.Ping(context.TODO(), nil)
	if err != nil {
		panic(err)
	}
	fmt.Println("Pinged successfully")
}

func DisconnectDB() {
	fmt.Println("Closing DB")
	if mongoClient != nil {
		mongoClient.Disconnect(context.TODO())
	}
	fmt.Println("Closed DB")
}

func getDB() (*mongo.Database, error) {
	if mongoClient == nil {
		return nil, fmt.Errorf("Not connected to DB")
	}
	return mongoClient.Database("url-shortner"), nil
}

func GetCollection(name string) (*mongo.Collection, error) {
	db, err := getDB()
	if err != nil {
		return nil, err
	}
	return db.Collection(name), nil
}
