package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

const SessionCollection = "sessions"

type Session struct {
	Id           bson.ObjectID `bson:"_id,omitempty"`
	SessionName  string        `bson:"session_name"`
	Description  string        `bson:"description"`
	Participants []Participant `bson:"participants"`
	Tags         []string      `bson:"tags"`
	CreatedAt    time.Time     `bson:"created_at"`
	EndedAt      *time.Time    `bson:"ended_at"`
}
