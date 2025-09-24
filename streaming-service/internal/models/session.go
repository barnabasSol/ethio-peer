package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

const SessionCollection = "sessions"

type Session struct {
	Id            bson.ObjectID `bson:"_id,omitempty"`
	OwnerId       string        `bson:"owner_id"`
	OwnerUsername string        `bson:"owner_username"`
	SessionName   string        `bson:"session_name"`
	Description   string        `bson:"description"`
	Participants  []Participant `bson:"participant_ids"`
	Tags          []string      `bson:"tags"`
	CreatedAt     time.Time     `bson:"created_at"`
	EndedAt       *time.Time    `bson:"ended_at"`
}
