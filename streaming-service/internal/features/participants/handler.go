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
	h.group.GET("/:session_id", h.GetParticipants)
	return h
}

func (h *Handler) Join(c echo.Context) error {
	username := c.Request().Header.Get("X-Claim-Username")
	user_id := c.Request().Header.Get("X-Claim-Sub")
	sid := c.Param("session_id")
	aa := c.QueryParam("as_anonymous")
	as_anonymous, err := strconv.ParseBool(aa)
	if err != nil {
		as_anonymous = false
	}
	name := c.QueryParam("name")
	profile_picture := c.QueryParam("profile_picture")

	req := Join{
		Name:           name,
		UserId:         user_id,
		Username:       username,
		ProfilePicture: profile_picture,
		SessionId:      sid,
		AsAnonymous:    as_anonymous,
	}
	if err := req.Validate(); err != nil {
		return c.JSON(
			http.StatusBadRequest,
			map[string]string{"error": err.Error()},
		)
	}

	res, err := h.s.Join(c.Request().Context(), req)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, res)
}

func (h *Handler) GetParticipants(c echo.Context) error {
	sid := c.Param("session_id")
	participants, err := h.s.GetParticipants(c.Request().Context(), sid)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, participants)
}
