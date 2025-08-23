package changepassword

import "github.com/labstack/echo/v4"

type Handler struct {
	group *echo.Group
	s     Service
}

func NewHandler(echo *echo.Echo) *Handler {
	h := &Handler{
		group: echo.Group("password"),
	}
	h.group.POST("/reset", h.ResetPassword)
	return h
}

func (h *Handler) ResetPassword(ctx echo.Context) error {
	return ctx.JSON(200, "alnjd")
}
