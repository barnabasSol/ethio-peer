package monitoring

import (
	"ep-streaming-service/internal/features/common/livekit"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	group  *echo.Group
	lk_cfg livekit.Config
	s      Service
}

func InitHandler(
	s Service,
	lk_cfg livekit.Config,
	group *echo.Group,
) *Handler {
	h := &Handler{
		group:  group,
		lk_cfg: lk_cfg,
		s:      s,
	}
	h.group.PATCH("/audio", h.Toggle)
	return h
}

func (h *Handler) Toggle(c echo.Context) error {
	username := c.Request().Header.Get("X-Claim-Username")
	var req ToggleAudio
	if err := c.Bind(&req); err != nil {
		return c.JSON(
			http.StatusBadRequest,
			map[string]string{"error": "invalid request body"},
		)
	}
	err := h.s.ToggleAudio(c.Request().Context(), req, username)
	if err != nil {
		return err
	}
	return c.JSON(
		http.StatusNoContent,
		map[string]string{"message": "muted successfully"},
	)
}
