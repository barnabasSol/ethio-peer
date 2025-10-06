package refreshtoken

import (
	"ep-auth-service/internal/features/common"
	"net/http"
	"os"
	"strconv"

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
	h.group.GET("", h.Refresh)
	return h
}
func (h *Handler) Refresh(ctx echo.Context) error {
	with_cookie := ctx.QueryParam("with_cookie")

	rt_cookie, err := ctx.Cookie("refresh_token")
	if err != nil || rt_cookie == nil || rt_cookie.Value == "" {
		return echo.NewHTTPError(
			http.StatusUnauthorized,
			"failed to refresh, no refresh token provided",
		)
	}
	result, err := h.s.Refresh(ctx.Request().Context(), Request{
		RefreshToken: rt_cookie.Value,
	})
	if err != nil {
		return err
	}

	if with_cookie != "" && with_cookie == "true" {
		if result.Data.AccessToken != "" && result.Data.RefreshToken != "" {
			expiry, err := strconv.Atoi(os.Getenv("JWT_EXPIRY_MINS"))
			if err != nil {
				return ctx.JSON(
					http.StatusInternalServerError,
					map[string]string{"error": "failed setting expiry"},
				)
			}

			cleared_at := common.ClearCookie("access_token")
			cleared_rt := common.ClearCookie("refresh_token")

			ctx.SetCookie(cleared_at)
			ctx.SetCookie(cleared_rt)

			common.SetCookie(
				"refresh_token",
				result.Data.RefreshToken,
				60*24*5,
			)
			atc := common.SetCookie(
				"access_token",
				result.Data.AccessToken,
				expiry,
			)
			rtc := common.SetCookie(
				"refresh_token",
				result.Data.RefreshToken,
				60*24*5,
			)
			ctx.SetCookie(atc)
			ctx.SetCookie(rtc)
		}
		return ctx.JSON(
			http.StatusOK,
			common.Response[Response]{
				Message: result.Message,
			},
		)
	}

	return ctx.JSON(http.StatusOK, result)
}
