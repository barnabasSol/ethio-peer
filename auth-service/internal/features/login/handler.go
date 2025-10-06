package login

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
	h.group.POST("/login", h.Login)
	return h
}

func (h *Handler) Login(ctx echo.Context) error {
	with_cookie := ctx.QueryParam("with_cookie")
	var req LoginRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(
			http.StatusBadRequest,
			map[string]string{"error": "invalid request"},
		)
	}

	if err := req.Validate(); err != nil {
		return err
	}

	result, err := h.s.LoginUser(ctx.Request().Context(), req)
	if err != nil {
		return err
	}
	if with_cookie != "" && with_cookie == "true" {
		if !result.Data.VerificationRequired {
			if result.Data.AccessToken != nil && result.Data.RefreshToken != nil {
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

				atc := common.SetCookie(
					"access_token",
					*result.Data.AccessToken,
					expiry,
				)
				rtc := common.SetCookie(
					"refresh_token",
					*result.Data.RefreshToken,
					60*24*7,
				)
				ctx.SetCookie(atc)
				ctx.SetCookie(rtc)
			}
			return ctx.JSON(
				http.StatusOK,
				common.Response[LoginResponse]{
					Message: result.Message,
					Data: LoginResponse{
						VerificationRequired: false,
					},
				},
			)
		}
	}
	return ctx.JSON(http.StatusOK, result)

}
