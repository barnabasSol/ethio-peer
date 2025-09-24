package participants

import (
	"ep-streaming-service/internal/features/common/livekit"
	"net/http"
	"strconv"

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
	h.group.POST("/join/:session_id", h.Join)
	return h
}

func (h *Handler) Join(c echo.Context) error {
	sid := c.Param("session_id")
	aa := c.QueryParam("as_anonymous")
	if aa == "" {
		aa = "false"
	}
	as_anonymous, err := strconv.ParseBool(aa)
	if err != nil {
		return c.JSON(
			http.StatusBadRequest,
			map[string]string{"error": "invalid param"},
		)
	}
	if sid == "" {
		return c.JSON(
			http.StatusBadRequest,
			map[string]string{"error": "please provide a session id"},
		)

	}

	user_id := c.Request().Header.Get("X-Claim-Sub")

	req := Join{
		UserId:      user_id,
		SessionId:   sid,
		AsAnonymous: as_anonymous,
	}
	res, err := h.s.Join(c.Request().Context(), req)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, res)
}
