package types

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	Username string             `json:"username" bson:"username"`
	ID       primitive.ObjectID `json:"id" bson:"_id,omitempty"`
}

type UserWithPassword struct {
	Password string `json:"password" bson:"password"`
	User     `bson:",inline"`
}
