package otp

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type OtpVerification struct {
	SessionId string `json:"session_id"`
	Code      string `json:"code"`
}

type PasswordOTPVerification struct {
	Code string `json:"code"`
}

type OtpSuccess struct {
	UserId       string  `json:"user_id"`
	AccessToken  *string `json:"access_token,omitempty"`
	RefreshToken *string `json:"refresh_token,omitempty"`
}

func (o OtpVerification) Validate() error {
	if o.SessionId == "" || o.Code == "" {
		return echo.NewHTTPError(
			http.StatusBadRequest,
			"invalid otp validation request",
		)
	}
	return nil
}
