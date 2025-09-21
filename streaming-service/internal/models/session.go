package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

const SessionCollection = "sessions"

type Session struct {
	Id             bson.ObjectID `bson:"_id,omitempty"`
	OwnerId        string        `json:"owner_id"`
	OwnerUsername  string        `json:"owner_username"`
	SessionName    string        `bson:"session_name"`
	Description    string        `json:"description"`
	ParticipantIds []string      `bson:"participant_ids"`
	IsLive         bool          `bson:"is_live"`
	CreatedAt      time.Time     `bson:"created_at"`
	EndedAt        time.Time     `bson:"updated_at"`
}
