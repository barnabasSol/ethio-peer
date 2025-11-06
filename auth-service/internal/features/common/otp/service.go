package otp

import (
	"context"
	"ep-auth-service/internal/features/common/cache"
	"ep-auth-service/internal/features/common/jwt"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type Service interface {
	VerifyAuthOTP(
		ctx context.Context,
		ov OtpVerification,
	) (*OtpSuccess, error)
	VerifyPasswordResetOTP(
		ctx context.Context,
		ov PasswordOTPVerification,
	) error
}

type service struct {
	repo Repository
	t    jwt.Generator
	c    *cache.Redis
	m    *OTPManager
}

func NewService(
	m *OTPManager,
	repo Repository,
	c *cache.Redis,
	t jwt.Generator,
) Service {
	return &service{
		repo: repo,
		t:    t,
		m:    m,
		c:    c,
	}
}

func (s *service) VerifyPasswordResetOTP(
	ctx context.Context,
	ov PasswordOTPVerification,
) error {
	key := fmt.Sprintf("psw_reset:%s", ov.Code)

	var institute_email string
	err := s.c.Get(ctx, key, &institute_email)
	if institute_email == "" || err != nil {
		return echo.NewHTTPError(
			http.StatusNotFound,
			"failed fetch identity, invalid otp",
		)
	}
	key = fmt.Sprintf("psw_reset_whitelist:%s", institute_email)
	err = s.c.SetWithTTL(
		ctx,
		key,
		institute_email,
		time.Minute*10,
	)
	if err != nil {
		return echo.NewHTTPError(
			http.StatusInternalServerError,
			"failed to set whitelist email",
		)
	}
	return nil
}
func (s *service) VerifyAuthOTP(
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

	refresh, err := s.t.GenerateRefreshTokenJWT(user.Id.Hex())
	if err != nil {
		log.Println(err)
		return nil, echo.NewHTTPError(
			http.StatusInternalServerError,
			"failed to generate refresh token",
		)
	}
	return &OtpSuccess{
		UserId:       user.Id.Hex(),
		AccessToken:  &token,
		RefreshToken: &refresh,
	}, nil
}
