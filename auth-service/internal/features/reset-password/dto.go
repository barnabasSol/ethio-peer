package resetpassword

type VerifyRequest struct {
	InstituteEmail string `json:"institute_email"`
}

type ChangePasswordRequest struct {
	InstituteEmail string `json:"institute_email"`
	NewPassword    string `json:"new_password"`
}
