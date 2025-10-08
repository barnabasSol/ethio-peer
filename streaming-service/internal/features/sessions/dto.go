package sessions

import "time"

type Create struct {
	OwnerProfilePic string   `json:"peer_profile_pic"`
	OwnerName       string   `json:"owner_name"`
	Name            string   `json:"name"`
	Topic           Topic    `json:"topic"`
	Description     string   `json:"description"`
	Tags            []string `json:"tags"`
}

type Update struct {
	SessionId   string  `json:"session_id"`
	SessionName *string `json:"session_name"`
	Description *string `json:"description"`
	Topic       *Topic  `json:"topic"`
	IsEnded     *bool   `json:"ended_at"`
}

type Topic struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type CreateResponse struct {
	RoomId string `json:"room_id"`
}

type Session struct {
	Name         string        `json:"name"`
	Owner        Owner         `json:"owner"`
	Duration     time.Duration `json:"duration"`
	Description  string        `json:"description"`
	CreatedAt    time.Time     `bson:"created_at"`
	EndedAt      *time.Time    `bson:"ended_at"`
	IsLive       bool          `json:"is_live"`
	Participants []Participant `json:"participants"`
}

type Owner struct {
	Name           string `json:"name"`
	Username       string `json:"username"`
	ProfilePicture string `json:"profile_picture"`
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
