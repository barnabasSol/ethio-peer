package login

type LoginRequest struct {
	Username       *string `json:"username,omitempty"`
	InstituteEmail *string `json:"institute_email,omitempty"`
	Email          *string `json:"email,omitempty"`
	Password       string  `json:"password"`
}

func (r *LoginRequest) Validate() error {
	count := 0
	if r.Username != nil {
		count++
	}
	if r.Email != nil {
		count++
	}
	if r.InstituteEmail != nil {
		count++
	}

	if count == 0 || count > 1 {
		return ErrInvalidCredential
	}
	return nil
}

type LoginResponse struct {
	VerificationRequired bool    `json:"verification_required"`
	OtpSessionId         *string `json:"otp_session_id,omitempty"`
	UserId               *string `json:"user_id,omitempty"`
	AccessToken          *string `json:"access_token,omitempty"`
	RefreshToken         *string `json:"refresh_token,omitempty"`
}
