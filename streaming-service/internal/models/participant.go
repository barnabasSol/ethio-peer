package models

import "time"

type Participant struct {
	UserId      string    `bson:"user_id"`
	IsAnonymous bool      `bson:"is_anonymous"`
	FlagStatus  int       `bson:"flag_status"`
	CreatedAt   time.Time `bson:"created_at"`
	UpdatedAt   time.Time `bson:"updated_at"`
}
