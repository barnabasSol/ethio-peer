package resetpassword

type VerifyRequest struct {
	InstituteEmail string `json:"institute_email"`
}

type ChangePasswordRequest struct {
	SessionId   string `json:"session_id"`
	NewPassword string `json:"new_password"`
}
