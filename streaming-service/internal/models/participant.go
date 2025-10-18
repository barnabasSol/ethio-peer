package models

import "time"

type Participant struct {
	UserId         string    `bson:"user_id"`
	Username       string    `bson:"username"`
	Name           string    `bson:"name"`
	ProfilePicture string    `bson:"profile_picture"`
	IsAnonymous    bool      `bson:"is_anonymous"`
	FlagStatus     int       `bson:"flag_status"`
	IsMuted        bool      `bson:"is_muted"`
	CreatedAt      time.Time `bson:"created_at"`
	UpdatedAt      time.Time `bson:"updated_at"`
}

type Owner struct {
	Username       string `bson:"username"`
	UserId         string `bson:"user_id"`
	Name           string `bson:"name"`
	ProfilePicture string `bson:"profile_picture"`
}
