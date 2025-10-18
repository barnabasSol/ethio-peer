package user

import "time"

type CurrentUser struct {
	UserId         string    `json:"user_id"`
	Username       string    `json:"username"`
	Name           string    `json:"name"`
	InstituteEmail string    `json:"institute_email"`
	Email          string    `json:"email,omitempty"`
	OverallScore   string    `json:"overall_score"`
	ProfilePhoto   string    `json:"profile_photo"`
	OnlineStatus   bool      `json:"online_status"`
	Bio            string    `json:"bio"`
	Roles          []string  `json:"roles"`
	Interests      []string  `json:"interests"`
	CreatedAt      time.Time `json:"created_at"`
}
