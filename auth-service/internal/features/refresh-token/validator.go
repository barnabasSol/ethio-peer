package refreshtoken

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (r *Request) Validate() error {
	if r.RefreshToken == "" {
		return echo.NewHTTPError(
			http.StatusUnauthorized,
			"failed to refresh",
		)
	}
	if r.UserId == "" {
		return echo.NewHTTPError(
			http.StatusUnauthorized,
			"invalid user",
		)
	}
	return nil
}
