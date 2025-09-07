package otp

import (
	"ep-auth-service/internal/features/shared"
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
		if status_code, found := shared.Errors[err]; found {
			return ctx.JSON(
				status_code,
				map[string]string{"error": err.Error()},
			)
		}
		return ctx.JSON(
			http.StatusInternalServerError,
			map[string]string{"error": "unexpected error"},
		)
	}

	return ctx.JSON(http.StatusOK, result)
}
