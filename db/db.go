package db

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoClient *mongo.Client

func ConnectDB() {
	godotenv.Load()
	MONGO_URI := os.Getenv("MONGO_URI")
	configOptions := options.Client().ApplyURI(MONGO_URI)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var err error
	mongoClient, err = mongo.Connect(ctx, configOptions)
	if err != nil {
		panic(err)
	}

	select {
	case <-ctx.Done():
		if ctx.Err() == context.DeadlineExceeded {
			fmt.Println("Connection attempt timed out")
			return
		}
	default:
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
