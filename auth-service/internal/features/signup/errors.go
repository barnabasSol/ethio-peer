package signup

import "errors"

var (
	ErrEmailAlreadyRegistered = errors.New("email is already registered")
	ErrInvalidSignupData      = errors.New("invalid signup data")
	ErrFailedToCreateUser     = errors.New("failed to create user")
	ErrFailedToSendOTP        = errors.New("failed to send OTP")
)

var SignupErrors = map[error]int{
	ErrEmailAlreadyRegistered: 409,
	ErrInvalidSignupData:      400,
	ErrFailedToCreateUser:     500,
	ErrFailedToSendOTP:        500,
}
