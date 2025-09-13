package broker

import "encoding/json"

type Message struct {
	Exchange string
	Topic    string          `json:"topic"`
	Data     json.RawMessage `json:"data"`
}

type PeerPayload struct {
	UserId    string   `json:"user_id"`
	Interests []string `json:"interests"`
	Bio       string   `json:"bio"`
}
