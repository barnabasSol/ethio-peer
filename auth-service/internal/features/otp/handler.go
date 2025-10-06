package otp

import (
	"ep-auth-service/internal/features/common"
	"net/http"

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
	h.group.POST("/verify", h.VerifyOTP)
	return h
}

func (h *Handler) VerifyOTP(ctx echo.Context) error {
	with_cookie := ctx.QueryParam("with_cookie")
	var req OtpVerification
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(
			http.StatusBadRequest,
			map[string]string{"error": "invalid request"},
		)
	}
	err := req.Validate()
	if err != nil {
		return err
	}

	result, err := h.s.VerifyOTP(ctx.Request().Context(), req)

	if err != nil {
		return err
	}

	if with_cookie != "" && with_cookie == "true" {

		cleared_at := common.ClearCookie("access_token")
		cleared_rt := common.ClearCookie("refresh_token")

		ctx.SetCookie(cleared_at)
		ctx.SetCookie(cleared_rt)

		atc := common.SetCookie("access_token", *result.AccessToken, 15)
		rtc := common.SetCookie("refresh_token", *result.RefreshToken, 60*24*7)

		ctx.SetCookie(atc)
		ctx.SetCookie(rtc)

		return ctx.JSON(
			http.StatusOK,
			common.Response[OtpSuccess]{
				Message: "login success",
				Data: OtpSuccess{
					UserId: result.UserId,
				},
			},
		)
	}

	return ctx.JSON(http.StatusOK, result)
}
