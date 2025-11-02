package logout

import (
	"github.com/labstack/echo/v4"
)

type Handler struct {
	group *echo.Group
	s     Service
}

func InitHandler(
	s Service,
	group *echo.Group,
) *Handler {
	h := &Handler{
		group: group,
		s:     s,
	}
	h.group.POST("/logout", h.Logout)
	return h
}

func (h *Handler) Logout(ctx echo.Context) error {
	return nil
}
