package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type PerformanceScore struct {
	UserId    bson.ObjectID `bson:"_id,omitempty"`
	Score     byte          `bson:"score"`
	CreatedAt time.Time     `bson:"created_at"`
}
