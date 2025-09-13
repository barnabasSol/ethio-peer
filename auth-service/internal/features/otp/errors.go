package otp

import "errors"

var (
	ErrInvalidOTP       = errors.New("invalid code")
	ErrIncorrectOTP     = errors.New("incorrect otp")
	ErrPendingOTP       = errors.New("you have a pending otp, please try again later")
	ErrFailedToGenOTP   = errors.New("failed to generate otp")
	ErrMissingOtpFields = errors.New("please send the required fields")
)

var OtpErrors = map[error]int{
	ErrInvalidOTP:       401,
	ErrMissingOtpFields: 400,
	ErrIncorrectOTP:     401,
	ErrPendingOTP:       401,
	ErrFailedToGenOTP:   500,
}
