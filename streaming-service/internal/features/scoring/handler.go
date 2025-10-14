package scoring

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
	lk_cfg livekit.Config,
	group *echo.Group,
	s Service,
) *Handler {
	h := &Handler{
		group:  group,
		lk_cfg: lk_cfg,
		s:      s,
	}
	h.group.POST("score", h.ScoreSession)
	return h
}

func (h *Handler) ScoreSession(c echo.Context) error {
	user_id := c.Request().Header.Get("X-Claim-Sub")
	var req Score
	if err := c.Bind(&req); err != nil {
		return c.JSON(
			http.StatusBadRequest,
			map[string]string{"error": "invalid request"},
		)
	}
	err := h.s.ScoreSession(c.Request().Context(), req, user_id)
	if err != nil {
		return err
	}
	return c.JSON(
		http.StatusBadRequest,
		map[string]string{"message": "rating accepted, thank you for your feedback"},
	)
}
