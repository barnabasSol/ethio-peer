package otp

import (
	"ep-auth-service/internal/features/common"
	"log"
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
		return ctx.JSON(
			http.StatusBadRequest,
			map[string]string{"error": err.Error()},
		)
	}
	result, err := h.s.VerifyOTP(ctx.Request().Context(), req)

	if err != nil {
		if status_code, found := OtpErrors[err]; found {
			return ctx.JSON(
				status_code,
				map[string]string{"error": err.Error()},
			)
		}
		if status_code, found := common.Errors[err]; found {
			return ctx.JSON(
				status_code,
				map[string]string{"error": err.Error()},
			)
		}
		log.Println(err)
		return ctx.JSON(
			http.StatusInternalServerError,
			map[string]string{"error": "unexpected error"},
		)
	}

	if with_cookie != "" && with_cookie == "true" {
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
