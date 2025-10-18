package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

const SessionCollection = "sessions"

type Session struct {
	Id            bson.ObjectID  `bson:"_id,omitempty"`
	SessionName   string         `bson:"session_name"`
	Description   string         `bson:"description"`
	Owner         Owner          `bson:"owner"`
	Topic         Topic          `bson:"topic"`
	ComputedScore string         `bson:"computed_score"`
	Scores        []SessionScore `bson:"scores"`
	Participants  []Participant  `bson:"participants"`
	Tags          []string       `bson:"tags"`
	CreatedAt     time.Time      `bson:"created_at"`
	UpdatedAt     time.Time      `bson:"updated_at"`
	StartsAt      *time.Time     `bson:"starts_at"`
	EndedAt       *time.Time     `bson:"ended_at"`
}

type Topic struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}
