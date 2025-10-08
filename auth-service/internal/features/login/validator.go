package login

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (r *LoginRequest) Validate() error {
	count := 0
	if r.Username != nil {
		count++
	}
	if r.Email != nil {
		count++
	}
	if r.InstituteEmail != nil {
		count++
	}
	if count == 0 || count > 1 {
		return echo.NewHTTPError(
			http.StatusBadRequest,
			"invalid credential",
		)
	}
	return nil
}
