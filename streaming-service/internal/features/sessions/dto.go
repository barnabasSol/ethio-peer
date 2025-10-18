package sessions

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Create struct {
	OwnerProfilePic string     `json:"peer_profile_pic"`
	OwnerName       string     `json:"owner_name"`
	Name            string     `json:"name"`
	Topic           Topic      `json:"topic"`
	Description     string     `json:"description"`
	StartsAt        *time.Time `json:"starts_at"`
	Tags            []string   `json:"tags"`
}

type Topic struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type Update struct {
	SessionId   string  `json:"session_id"`
	SessionName *string `json:"session_name"`
	Description *string `json:"description"`
	Topic       *Topic  `json:"topic"`
	IsEnded     *bool   `json:"is_ended"`
}

type CreateResponse struct {
	RoomId string `json:"room_id"`
}

type Session struct {
	Id           bson.ObjectID `bson:"_id" json:"-"`
	SessionId    string        `json:"session_id"`
	Name         string        `bson:"session_name" json:"name"`
	Owner        Owner         `bson:"owner" json:"owner"`
	Description  string        `bson:"description" json:"description"`
	StartsAt     time.Time     `bson:"starts_at" json:"starts_at"`
	EndedAt      *time.Time    `bson:"ended_at,omitempty" json:"ended_at,omitempty"`
	Participants []Participant `bson:"participants" json:"participants"`
	Duration     string        `json:"duration"`
	IsLive       bool          `json:"is_live"`
}

type Owner struct {
	Username       string `bson:"username"`
	Name           string `bson:"name"`
	ProfilePicture string `bson:"profile_picture"`
}

type Participant struct {
	Username       string `json:"username"`
	ProfilePicture string `json:"profile_picture"`
}

type Join struct {
	SessionId   string `json:"session_id"`
	UserId      string `json:"user_id"`
	IsAnonymous string `json:"is_anaonymous"`
}
