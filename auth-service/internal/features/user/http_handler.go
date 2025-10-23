package user

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type HttpHandler struct {
	s Service
	e *echo.Group
}

func InitHandler(e *echo.Group, s Service) {
	h := &HttpHandler{
		e: e,
		s: s,
	}
	h.e.GET("/user/count", h.Count)
}

func (h *HttpHandler) Count(c echo.Context) error {
	count, err := h.s.GetUserCount(c.Request().Context())
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, count)
}
