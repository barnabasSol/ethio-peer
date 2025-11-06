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

func (r *ChangePasswordRequest) Validate() error {
	if r.InstituteEmail == "" {
		return echo.NewHTTPError(
			http.StatusBadRequest,
			"please provide correct credential",
		)
	}
	if r.NewPassword == "" {
		return echo.NewHTTPError(
			http.StatusBadRequest,
			"please provide new password",
		)
	}
	return nil
}
