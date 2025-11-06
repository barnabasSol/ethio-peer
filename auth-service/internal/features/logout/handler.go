package logout

import (
	"ep-auth-service/internal/features/common"
	"net/http"

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

	cleared_at := common.ClearCookie("access_token")
	cleared_rt := common.ClearCookie("refresh_token")

	ctx.SetCookie(cleared_at)
	ctx.SetCookie(cleared_rt)

	return ctx.JSON(
		http.StatusNoContent,
		"successfully logged out",
	)
}
