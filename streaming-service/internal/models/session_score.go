package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type SessionScore struct {
	Id        bson.ObjectID `bson:"_id,omitempty"`
	UserId    string        `bson:"user_id"`
	Score     float32       `bson:"score"`
	Comment   string        `bson:"comment"`
	CreatedAt time.Time     `bson:"created_at"`
}
