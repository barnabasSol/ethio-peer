package login

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
	h.group.POST("/login", h.Login)
	return h
}

func (h *Handler) Login(ctx echo.Context) error {
	var req LoginRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(
			http.StatusBadRequest,
			map[string]string{"error": "invalid request"},
		)
	}
	if err := req.Validate(); err != nil {
		return ctx.JSON(
			LoginErrors[err],
			map[string]string{"error": err.Error()},
		)
	}

	result, err := h.s.LoginUser(ctx.Request().Context(), req)

	if err != nil {
		if status_code, found := LoginErrors[err]; found {
			return ctx.JSON(
				status_code,
				map[string]string{"error": err.Error()},
			)
		}
		return ctx.JSON(
			http.StatusInternalServerError,
			map[string]string{"error": "unexpected error"},
		)
	}

	return ctx.JSON(http.StatusOK, result)

}
