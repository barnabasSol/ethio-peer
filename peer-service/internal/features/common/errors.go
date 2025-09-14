package common

import "errors"

var (
	ErrPeerNotFound  = errors.New("peer not found")
	ErrInvalidPeerId = errors.New("invalid peer id")
)

var Errors = map[error]int{
	ErrPeerNotFound:  404,
	ErrInvalidPeerId: 400,
}
