package models

import "time"

type Participant struct {
	Username       string    `bson:"username"`
	Name           string    `bson:"name"`
	ProfilePicture string    `bson:"profile_picture"`
	IsAnonymous    bool      `bson:"is_anonymous"`
	IsOwner        bool      `bson:"is_owner"`
	FlagStatus     int       `bson:"flag_status"`
	CreatedAt      time.Time `bson:"created_at"`
	UpdatedAt      time.Time `bson:"updated_at"`
}
