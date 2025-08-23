package signup

type SignUpRequest struct {
	Name           string      `json:"name"`
	Username       string      `json:"username"`
	InstituteEmail string      `json:"institute_email"`
	Email          string      `json:"email"`
	Password       string      `json:"password"`
	Interests      *[]Interest `json:"interests"`
	Bio            *string     `json:"bio,omitempty"`
}

type Interest struct {
	Id    string `json:"id"`
	Topic string `json:"topic"`
}

type SignUpResponse struct {
	UserId       string `json:"user_id"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
