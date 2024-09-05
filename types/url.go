package types

import "go.mongodb.org/mongo-driver/bson/primitive"

type Url struct {
	URL    string             `json:"url" bson:"url"`
	ID     primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserId primitive.ObjectID `json:"userId" bson:"user"`
}

type UrlWithShortVersion struct {
	ShortenedUrl string `json:"shortenedUrl" bson:"shUrl"`
	Url          `bson:",inline"`
}

type UrlWithUser struct {
	User User
	Url  `bson:",inline"`
}
