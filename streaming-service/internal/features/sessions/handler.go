package sessions

import (
	"ep-streaming-service/internal/features/common/livekit"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/livekit/protocol/auth"
	"github.com/livekit/protocol/webhook"
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
	h.group.POST("", h.CreateSession)
	h.group.PATCH("", h.UpdateSession)
	h.group.POST("/livekit/webhook", h.HandleLiveKitWebhook)
	return h
}

func (h *Handler) HandleLiveKitWebhook(c echo.Context) error {
	authProvider := auth.NewSimpleKeyProvider(
		h.lk_cfg.ApiKey,
		h.lk_cfg.ApiSecret,
	)
	event, err := webhook.ReceiveWebhookEvent(
		c.Request(),
		authProvider,
	)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	if event.Event == "participant_joined" {
		participant := event.Participant
		room := event.Room
		fmt.Printf(
			"Participant %s joined room %s\n",
			participant.Identity,
			room.Name,
		)
	}

	return c.NoContent(http.StatusOK)
}
func (h *Handler) CreateSession(c echo.Context) error {
	user_id := c.Request().Header.Get("X-Claim-Sub")
	username := c.Request().Header.Get("X-Claim-Username")
	var req Create
	if err := c.Bind(&req); err != nil {
		return c.JSON(
			http.StatusBadRequest,
			map[string]string{"error": "invalid request"},
		)
	}
	if err := req.Validate(); err != nil {
		return err
	}
	res, err := h.s.CreteSession(
		c.Request().Context(),
		username,
		user_id,
		req,
	)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, res)
}

func (h *Handler) UpdateSession(c echo.Context) error {
	username := c.Request().Header.Get("X-Claim-Username")
	var req Update
	if err := c.Bind(&req); err != nil {
		return c.JSON(
			http.StatusBadRequest,
			map[string]string{"error": "invalid request"},
		)
	}
	if err := req.Validate(); err != nil {
		return c.JSON(
			http.StatusBadRequest,
			map[string]string{"error": err.Error()},
		)
	}

	err := h.s.UpdateSession(
		c.Request().Context(),
		req,
		username,
	)

	if err != nil {
		return err
	}

	return c.JSON(
		http.StatusOK,
		map[string]string{"result": "successfully updated"},
	)
}
