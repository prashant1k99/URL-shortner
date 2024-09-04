package db

import (
	"context"

	"github.com/prashant1k99/URL-Shortner/types"
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
