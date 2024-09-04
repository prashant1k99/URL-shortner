package types

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Analytics struct {
	ID     primitive.ObjectID `json:"id" bson:"_id"`
	UrlId  primitive.ObjectID `json:"urlId" bson:"urlId"`
	IP     string             `json:"ip" bson:"ip"`
	AtTime time.Time          `json:"at" bson:"at"`
}
