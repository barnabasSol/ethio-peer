package sessions

import "time"

type Create struct {
	Name        string   `json:"name"`
	Subject     string   `json:"subject"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
}
type CreateResponse struct {
	// Token  string `json:"token"`
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
