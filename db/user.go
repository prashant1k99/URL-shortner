package db

import (
	"context"

	"github.com/prashant1k99/URL-Shortner/types"
	"github.com/prashant1k99/URL-Shortner/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetUserById(id primitive.ObjectID) (types.UserWithPassword, error) {
	userCollection, err := GetCollection("users")
	if err != nil {
		return types.UserWithPassword{}, err
	}
	var user types.UserWithPassword
	err = userCollection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&user)
	if err != nil {
		return types.UserWithPassword{}, err
	}
	return user, nil
}

func GetUserByUsername(username string) (types.UserWithPassword, error) {
	userCollection, err := GetCollection("users")
	if err != nil {
		return types.UserWithPassword{}, err
	}
	var user types.UserWithPassword
	err = userCollection.FindOne(context.TODO(), bson.M{"username": username}).Decode(&user)
	if err != nil {
		return types.UserWithPassword{}, err
	}
	return user, nil
}

func CreateUser(user *types.UserWithPassword) (primitive.ObjectID, error) {
	userCollection, err := GetCollection("users")
	if err != nil {
		return primitive.ObjectID{}, err
	}
	encryptedPass, err := utils.EncryptPassword(user.Password)
	if err != nil {
		return primitive.ObjectID{}, err
	}
	result, err := userCollection.InsertOne(context.TODO(), bson.M{
		"username": user.Username,
		"password": encryptedPass,
	})
	if err != nil {
		return primitive.ObjectID{}, err
	}
	return result.InsertedID.(primitive.ObjectID), nil
}
