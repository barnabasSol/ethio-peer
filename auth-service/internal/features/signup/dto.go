package signup

type SignUpRequest struct {
	Name           string    `json:"name"`
	Username       string    `json:"username"`
	InstituteEmail string    `json:"institute_email"`
	Email          string    `json:"email"`
	Password       string    `json:"password"`
	Interests      *[]string `json:"interests,omitempty"`
	Bio            *string   `json:"bio,omitempty"`
}

type SignUpResponse struct {
	VerificationRequired bool    `json:"verification_required"`
	OtpSessionId         *string `json:"otp_session_id,omitempty"`
	UserId               *string `json:"user_id,omitempty"`
	AccessToken          *string `json:"access_token,omitempty"`
	RefreshToken         *string `json:"refresh_token,omitempty"`
}
