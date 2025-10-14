package resetpassword

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (r *VerifyRequest) Validate() error {
	if r.InstituteEmail != "" {
		return echo.NewHTTPError(
			http.StatusBadRequest,
			"please provide correct credential",
		)
	}
	return nil
}
