package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

const TokenCollection = "refresh_tokens"

type RefreshToken struct {
	Id           bson.ObjectID `bson:"_id,omitempty"`
	UserId       bson.ObjectID `bson:"user_id"`
	RefreshToken string        `bson:"refresh_token"`
	CreatedAt    time.Time     `bson:"created_at"`
	UpdatedAt    time.Time     `bson:"updated_at"`
}
