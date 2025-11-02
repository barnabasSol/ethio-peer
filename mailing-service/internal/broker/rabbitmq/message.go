package broker

import "encoding/json"

type Message struct {
	Exchange string          `json:"exchange"`
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

type PasswordResetPayload struct {
	Email string `json:""`
}
