package otp

import "errors"

var (
	ErrInvalidOTP       = errors.New("invalid code")
	ErrIncorrectOTP     = errors.New("incorrect otp")
	ErrMissingOtpFields = errors.New("please send the required fields")
)

var OtpErrors = map[error]int{
	ErrInvalidOTP:       401,
	ErrMissingOtpFields: 400,
	ErrIncorrectOTP:     401,
}
