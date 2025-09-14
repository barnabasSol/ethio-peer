package peer

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	s     Service
	group *echo.Group
}

func InitHandler(s Service, group *echo.Group) *Handler {
	h := &Handler{
		group: group,
		s:     s,
	}
	h.group.GET("/peer/:user_id", h.GetPeer)
	return h
}

func (h *Handler) GetPeer(ctx echo.Context) error {
	user_id := ctx.Param("user_id")
	peer, err := h.s.GetPeer(ctx.Request().Context(), user_id)
	if err != nil {
		log.Println(err)
		return ctx.JSON(
			http.StatusInternalServerError,
			map[string]string{
				"error": err.Error(),
			},
		)
	}
	return ctx.JSON(http.StatusOK, peer)

}
