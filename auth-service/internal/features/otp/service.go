package otp

import (
	"context"
	"ep-auth-service/internal/features/common/jwt"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Service interface {
	VerifyOTP(
		ctx context.Context,
		ov OtpVerification,
	) (*OtpSuccess, error)
}

type service struct {
	repo Repository
	t    jwt.Generator
	m    *OTPManager
}

func NewService(
	m *OTPManager,
	repo Repository,
	t jwt.Generator,
) Service {
	return &service{
		repo: repo,
		t:    t,
		m:    m,
	}
}

func (s *service) VerifyOTP(
	ctx context.Context,
	ov OtpVerification,
) (*OtpSuccess, error) {
	s.m.mu.RLock()
	v, found := s.m.collection[ov.SessionId]
	s.m.mu.RUnlock()
	if !found {
		return nil, echo.NewHTTPError(
			http.StatusBadRequest,
			"invalid otp session",
		)
	}
	if v.Value != ov.Code {
		return nil, echo.NewHTTPError(
			http.StatusUnauthorized,
			"incorrect code",
		)
	}

	s.m.removeOTP(ov.SessionId)

	user, err := s.repo.GetUserById(ctx, v.UserId)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	err = s.repo.UpdateUser(ctx, user.Id, true, true)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	token, err := s.t.GenerateAccessToken(*user)
	if err != nil {
		log.Println(err)
		return nil, echo.NewHTTPError(
			http.StatusInternalServerError,
			"failed to authenticate, try again later",
		)
	}

	refresh, err := s.t.GenerateRefreshToken(32)
	if err != nil {
		log.Println(err)
		return nil, echo.NewHTTPError(
			http.StatusInternalServerError,
			"failed to authenticate, try again later",
		)
	}
	return &OtpSuccess{
		UserId:       user.Id.Hex(),
		AccessToken:  &token,
		RefreshToken: &refresh,
	}, nil
}
