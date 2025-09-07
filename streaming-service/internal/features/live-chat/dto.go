package livechat

import "time"

type MessageCreate struct {
	From        string `json:"from"`
	ToRoom      string `json:"to_room"`
	Message     string `json:"message"`
	IsAnonymous bool   `json:"is_anonymous"`
}

type MessageResponse struct {
	From      string    `json:"from"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}
