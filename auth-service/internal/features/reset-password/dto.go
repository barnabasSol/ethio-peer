package resetpassword

type VerifyRequest struct {
	Username       *string `json:"username,omitempty"`
	Email          *string `json:"email,omitempty"`
	InstituteEmail *string `json:"institute_email,omitempty"`
}
type ChangePasswordRequest struct {
	SessionId   string `json:"session_id"`
	NewPassword string `json:"new_password"`
}
