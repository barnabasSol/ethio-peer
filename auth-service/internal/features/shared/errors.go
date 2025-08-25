package shared

import "errors"

var Errors = map[int]error{
	404: errors.New("not found"),
	409: errors.New("validation error"),
}
