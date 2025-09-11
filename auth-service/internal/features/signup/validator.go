package signup

import (
	"errors"
	"regexp"
	"strings"
)

func (s *SignUpRequest) Validate() error {
	if strings.TrimSpace(s.Name) == "" {
		return errors.New("name is required")
	}
	if strings.TrimSpace(s.Username) == "" {
		return errors.New("username is required")
	}
	if strings.TrimSpace(s.InstituteEmail) == "" {
		return errors.New("institute_email is required")
	}
	if !isValidEmail(s.InstituteEmail) {
		return errors.New("institute_email is not a valid email")
	}
	if !strings.HasSuffix(s.InstituteEmail, "@hilcoeschool.com") {
		return errors.New("institute_email must end with @hilcoeschool.com")
	}
	if strings.TrimSpace(s.Email) == "" {
		return errors.New("email is required")
	}
	if !isValidEmail(s.Email) {
		return errors.New("email is not a valid email")
	}
	if len(s.Password) < 6 {
		return errors.New("password must be at least 6 characters")
	}
	if s.Interests != nil && len(*s.Interests) == 0 {
		return errors.New("interests must not be empty if provided")
	}
	return nil
}

func isValidEmail(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}
