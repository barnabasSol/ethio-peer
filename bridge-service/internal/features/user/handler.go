package user

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
	h.group.GET("/me", h.GetUser)
	return h
}

func (h *Handler) GetUser(ctx echo.Context) error {
	user_id := ctx.Request().Header.Get("X-Claim-Sub")
	peer, err := h.s.GetCurrentUser(ctx.Request().Context(), user_id)
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
