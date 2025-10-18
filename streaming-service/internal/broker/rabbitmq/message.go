package broker

import "encoding/json"

type Message struct {
	Exchange string
	Topic    string          `json:"topic"`
	Data     json.RawMessage `json:"data"`
}

type OtpPayload struct {
	Email string `json:"email"`
	OTP   string `json:"otp"`
}

type WelcomePayload struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

type ScorePayload struct {
	UserId string  `json:"user_id"`
	Score  float32 `json:"score"`
}

type NewSessionPayload struct {
	OwnerId   string `json:"OwnerId"`
	SessionId string `json:"SessionId"`
	UserName  string `json:"UserName"`
	TopicId   string `json:"TopicId"`
}

type NewParticipantPayload struct {
	SessionId string `json:"SessionId"`
	MemberId  string `json:"MemberId"`
}
