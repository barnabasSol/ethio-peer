package profile

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	group *echo.Group
	s     Service
}

func InitHandler(s Service, group *echo.Group) *Handler {
	h := &Handler{
		group: group,
		s:     s,
	}
	h.group.PUT("/photo", h.UpdateProfilePicture)
	h.group.DELETE("/photo", h.DeleteProfilePicture)
	return h
}

func (h *Handler) UpdateProfilePicture(c echo.Context) error {
	sub := c.Request().Header.Get("X-Claim-Sub")
	if sub == "" {
		return echo.NewHTTPError(
			http.StatusUnauthorized,
			"missing X-Claim-Sub header",
		)
	}
	presignedURL, err := h.s.UpdateProfilePicture(
		c.Request().Context(),
		sub,
	)
	if err != nil {
		return err
	}
	return c.JSON(
		http.StatusOK,
		map[string]string{
			"upload_url": presignedURL,
		},
	)
}

func (h *Handler) DeleteProfilePicture(c echo.Context) error {
	sub := c.Request().Header.Get("X-Claim-Sub")
	if sub == "" {
		return echo.NewHTTPError(
			http.StatusUnauthorized,
			"missing X-Claim-Sub header",
		)
	}
	err := h.s.DeleteProfilePicture(c.Request().Context(), sub)
	if err != nil {
		return err
	}
	return c.JSON(
		http.StatusOK,
		map[string]string{
			"message": "profile picture deleted successfully",
		},
	)
}
