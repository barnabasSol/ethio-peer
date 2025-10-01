package participants

import (
	"errors"
	"strings"
)

func (j Join) Validate() error {
	if strings.TrimSpace(j.Name) == "" {
		return errors.New("name is required")
	}
	if strings.TrimSpace(j.Username) == "" {
		return errors.New("username is required")
	}
	if strings.TrimSpace(j.SessionId) == "" {
		return errors.New("session_id is required")
	}
	if j.ProfilePicture != "" && !strings.HasPrefix(j.ProfilePicture, "http") {
		return errors.New("profile_picture must be a valid URL")
	}
	return nil
}
