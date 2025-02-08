package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type DiscordUser struct {
	BId      primitive.ObjectID `json:"_id" bson:"_id"` //nolint
	AuthorId string             `json:"author_id" bson:"author_id"`
	Lang     string             `json:"lang" bson:"lang"`
}
