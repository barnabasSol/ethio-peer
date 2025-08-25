package login

import "errors"

var (
	ErrIncorrectPassword = errors.New("incorrect password")
	ErrInvalidCredential = errors.New("invalid credential")
	ErrUserDoesntExist   = errors.New("user doesnt exist")
)

var LoginErrors = map[error]int{
	ErrIncorrectPassword: 401,
	ErrInvalidCredential: 401,
	ErrUserDoesntExist:   404,
}
