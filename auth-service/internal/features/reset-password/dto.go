package changepassword

type Request struct {
	UserId      string `json:"user_id"`
	OTP         string `json:"otp"`
	NewPassword string `json:"new_password"`
}
