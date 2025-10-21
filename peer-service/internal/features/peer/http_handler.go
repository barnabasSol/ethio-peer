package peer

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	e echo.Group
	s Service
}

func NewHandler(
	s Service,
	e echo.Group,
) *Handler {
	h := &Handler{
		s: s,
		e: e,
	}
	h.e.GET("", h.GetTopPeers)
	return h
}

func (h *Handler) GetTopPeers(c echo.Context) error {
	res, err := h.s.GetTopPeers(c.Request().Context())
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, res)

}
