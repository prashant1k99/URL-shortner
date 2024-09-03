package types

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID       primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Username string             `json:"username" bson:"username"`
}

type UserWithPassword struct {
	User     `bson:",inline"`
	Password string `json:"password" bson:"password"`
}
