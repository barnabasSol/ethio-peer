package signup

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
	h.group.POST("/signup", h.SignUpUser)
	return h
}

func (h *Handler) SignUpUser(ctx echo.Context) error {
	var req SignUpRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(
			http.StatusBadRequest,
			map[string]string{"error": "invalid request"},
		)
	}

	err := req.Validate()

	if err != nil {
		return ctx.JSON(
			http.StatusBadRequest,
			map[string]string{"error": err.Error()},
		)
	}

	result, err := h.s.SignUpUser(ctx.Request().Context(), req)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, result)
}
