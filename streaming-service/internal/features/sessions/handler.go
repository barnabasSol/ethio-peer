package sessions

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
	h.group.POST("", h.CreateSession)
	return h
}

func (h *Handler) CreateSession(c echo.Context) error {
	user_id := c.Request().Header.Get("X-Claim-Sub")
	username := c.Request().Header.Get("X-Claim-Username")
	_ = user_id
	var req Create
	if err := c.Bind(&req); err != nil {
		return c.JSON(
			http.StatusBadRequest,
			map[string]string{"error": "invalid request"},
		)
	}
	if err := req.Validate(); err != nil {
		return err
	}

	return nil
}
