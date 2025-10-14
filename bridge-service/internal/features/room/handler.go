package room

import "github.com/labstack/echo/v4"

type Handler struct {
	s     Service
	group *echo.Group
}

func InitHandler(
	s Service,
	group *echo.Group,
) *Handler {
	h := &Handler{
		group: group,
		s:     s,
	}
	h.group.GET(
		"/:room_id",
		h.GetRoom,
	)
	return h
}

func (h *Handler) GetRoom(c echo.Context) error {
	return nil
}
