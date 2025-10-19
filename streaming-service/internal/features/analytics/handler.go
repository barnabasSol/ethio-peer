package analytics

import (
	"ep-streaming-service/internal/features/common"
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
	group *echo.Group,
) *Handler {
	h := &Handler{
		group: group,
		s:     s,
	}
	h.group.GET("", h.GetSessionAnalytics)
	return h
}

func (h *Handler) GetSessionAnalytics(c echo.Context) error {
	filter := c.QueryParam("filter")
	res, err := h.s.GetSessionAnalytics(c.Request().Context(), filter)
	if err != nil {
		return err
	}
	result := common.Response[*SessionAnalytics]{
		Message: "success",
		Data:    res,
	}
	return c.JSON(http.StatusOK, result)
}
