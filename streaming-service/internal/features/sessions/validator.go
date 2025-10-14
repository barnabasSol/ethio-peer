package sessions

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (r Create) Validate() error {
	if r.Name == "" {
		return echo.NewHTTPError(
			http.StatusBadRequest,
			"please provide a name for the session",
		)
	}

	if len(r.Name) > 100 {
		return echo.NewHTTPError(
			http.StatusBadRequest,
			"name must be at most 100 characters",
		)
	}

	if len(r.Description) > 500 {
		return echo.NewHTTPError(
			http.StatusBadRequest,
			"description must be at most 500 characters",
		)
	}

	for i, tag := range r.Tags {
		if tag == "" {
			return echo.NewHTTPError(
				http.StatusBadRequest,
				fmt.Sprintf("tag at position %d cannot be empty", i),
			)
		}
		if len(tag) > 30 {
			return echo.NewHTTPError(
				http.StatusBadRequest,
				fmt.Sprintf("tag '%s' is too long (max 30 chars)", tag),
			)
		}
	}

	return nil
}

func (r Update) Validate() error {
	if r.SessionId == "" {
		return errors.New("please provide the session id")
	}
	count := 0
	if r.Description == nil {
		count++
	}
	if r.IsEnded == nil {
		count++

	}
	if r.SessionName == nil {
		count++
	}
	if r.Topic == nil {
		count++
	}
	if count == 4 {
		return errors.New("please provide neccessary fields")
	}

	return nil
}

var ValidFilters = map[string]struct{}{
	"upcoming":  {},
	"ongoing":   {},
	"concluded": {},
	"":          {},
}
