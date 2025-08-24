package login

import "errors"

type LoginRequest struct {
	Username       *string `json:"username,omitempty"`
	Email          *string `json:"email,omitempty"`
	InstituteEmail *string `json:"institute_email,omitempty"`
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

	if count == 0 {
		return errors.New("no credential provided, one must be provided")
	}
	if count > 1 {
		return errors.New("only one credential must be provided")
	}
	return nil
}

type LoginReponse struct {
	UserId       string `json:"user_id"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
