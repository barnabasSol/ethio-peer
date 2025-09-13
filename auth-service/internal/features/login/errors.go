package login

import "errors"

var (
	ErrIncorrectPassword    = errors.New("incorrect password")
	ErrInvalidCredential    = errors.New("invalid credential")
	ErrUserDoesntExist      = errors.New("user doesnt exist")
	ErrFailedToAuthenticate = errors.New("failed to authenticate please try again later")
)

var LoginErrors = map[error]int{
	ErrIncorrectPassword:    401,
	ErrInvalidCredential:    401,
	ErrUserDoesntExist:      404,
	ErrFailedToAuthenticate: 500,
}
