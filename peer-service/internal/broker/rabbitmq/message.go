package broker

import "encoding/json"

type Message struct {
	Exchange string
	Topic    string          `json:"topic"`
	Data     json.RawMessage `json:"data"`
}

type PeerPayload struct {
	UserId    string     `json:"user_id"`
	Interests []Interest `json:"interests"`
	Bio       string     `json:"bio"`
}

type Interest struct {
	Id    string `json:"id"`
	Topic string `json:"topic"`
}
