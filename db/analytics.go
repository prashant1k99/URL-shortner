package db

import (
	"context"

	"github.com/prashant1k99/URL-Shortner/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AddAnalytics(analytic types.Analytics) (primitive.ObjectID, error) {
	analyticsCollection, err := GetCollection("analytics")
	if err != nil {
		return primitive.ObjectID{}, err
	}
	result, err := analyticsCollection.InsertOne(context.TODO(), analytic)
	if err != nil {
		return primitive.ObjectID{}, err
	}
	return result.InsertedID.(primitive.ObjectID), nil
}

func GetUrlVisitCounts(urlId primitive.ObjectID) (int64, error) {
	analyticsCollection, err := GetCollection("analytics")
	if err != nil {
		return 0, err
	}
	count, err := analyticsCollection.CountDocuments(context.TODO(), bson.M{"urlId": urlId})
	if err != nil {
		return 0, err
	}
	return count, nil
}
