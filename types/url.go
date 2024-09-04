package types

import "go.mongodb.org/mongo-driver/bson/primitive"

type Url struct {
	ID     primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	URL    string             `json:"url" bson:"url"`
	UserId primitive.ObjectID `json:"userId" bson:"user"`
}
type UrlWithShortVersion struct {
	Url          `bson:",inline"`
	ShortenedUrl string `json:"shortenedUrl" bson:"shUrl"`
}
type UrlWithUser struct {
	Url  `bson:",inline"`
	User User
}
