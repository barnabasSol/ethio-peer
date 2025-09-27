package refreshtoken

type Request struct {
	UserId       string `json:"user_id"`
	RefreshToken string `json:"refresh_token"`
}

type Response struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
