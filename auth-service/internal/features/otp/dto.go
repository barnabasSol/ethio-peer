package otp

type OtpVerification struct {
	SessionId string `json:"session_id"`
	Code      string `json:"code"`
}

type OtpSuccess struct {
	UserId       string `json:"user_id,omitempty"`
	AccessToken  string `json:"access_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

func (o OtpVerification) Validate() error {
	if o.SessionId == "" || o.Code == "" {
		return ErrMissingOtpFields
	}
	return nil
}
