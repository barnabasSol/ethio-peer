package resetpassword

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	group *echo.Group
	s     Service
}

func InitHandler(echo *echo.Echo, s Service) *Handler {
	h := &Handler{
		group: echo.Group("password"),
		s:     s,
	}
	h.group.POST("/reset", h.ResetPassword)
	h.group.POST("/email/otp", h.SendResetPasswordEmail)
	return h
}

func (h *Handler) ResetPassword(ctx echo.Context) error {
	var req ChangePasswordRequest
	if err := ctx.Bind(&req); err != nil {
		return echo.NewHTTPError(
			http.StatusBadRequest,
			"invalid request",
		)
	}
	if err := req.Validate(); err != nil {
		return err
	}

	if err := h.s.ResetPassword(ctx.Request().Context(), req); err != nil {
		return err
	}

	return ctx.JSON(
		http.StatusNoContent,
		"successfully reset",
	)
}

func (h *Handler) SendResetPasswordEmail(ctx echo.Context) error {
	var req VerifyRequest
	if err := ctx.Bind(&req); err != nil {
		return echo.NewHTTPError(
			http.StatusBadRequest,
			"invalid request",
		)
	}
	err := h.s.VerifyCredentialAndSendOTP(ctx.Request().Context(), req)
	if err != nil {
		return err
	}
	return ctx.JSON(
		http.StatusNoContent,
		"check your email",
	)
}
