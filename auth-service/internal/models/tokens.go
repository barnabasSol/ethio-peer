package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

const TokenCollection = "tokens"

type Tokens struct {
	Id        bson.ObjectID `bson:"_id,omitempty"`
	UserId    bson.ObjectID `bson:"user_id"`
	CreatedAt time.Time     `bson:"created_at"`
	UpdatedAt time.Time     `bson:"updated_at"`
}
