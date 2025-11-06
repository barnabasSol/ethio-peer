package resetpassword

import (
	"context"
	"encoding/json"
	broker "ep-auth-service/internal/broker/rabbitmq"
	"ep-auth-service/internal/features/common/cache"
	"ep-auth-service/internal/features/common/otp"
	"ep-auth-service/internal/features/user"
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	VerifyCredentialAndSendOTP(context.Context, VerifyRequest) error
	ResetPassword(context.Context, ChangePasswordRequest) error
}

type service struct {
	s     *broker.RabbitMQ
	repo  Repository
	urepo user.Repository
	otpm  *otp.OTPManager
	cache *cache.Redis
	br    *broker.RabbitMQ
}

func NewService(
	repo Repository,
	cache *cache.Redis,
	br *broker.RabbitMQ,
	urepo user.Repository,
	otpm *otp.OTPManager,
) Service {
	return &service{
		repo:  repo,
		br:    br,
		cache: cache,
		urepo: urepo,
		otpm:  otp.NewOTPManager(context.Background()),
	}
}

func (s *service) VerifyCredentialAndSendOTP(
	ctx context.Context,
	req VerifyRequest,
) error {
	u, err := s.urepo.GetUserByInstitueEmail(ctx, req.InstituteEmail)

	if err != nil {
		return err
	}

	o, err := s.otpm.Generate(u.Id.Hex())
	if err != nil {
		return err
	}

	otp := broker.OtpPayload{
		Email: u.InstituteEmail,
		OTP:   o.Value,
	}

	otp_json, err := json.Marshal(otp)
	if err != nil {
		return echo.NewHTTPError(
			http.StatusInternalServerError,
			"failed to bind otp",
		)
	}
	key := fmt.Sprintf("psw_reset:%s", o.Value)
	err = s.cache.Set(ctx, key, u.InstituteEmail)
	if err != nil {
		return echo.NewHTTPError(
			http.StatusInternalServerError,
			"request for otp later",
		)
	}

	s.br.Publish(broker.Message{
		Exchange: "notification_exchange",
		Topic:    "email.otp",
		Data:     otp_json,
	})

	return nil
}

func (s *service) ResetPassword(
	ctx context.Context,
	req ChangePasswordRequest,
) error {
	key := fmt.Sprintf("psw_reset_whitelist:%s", req.InstituteEmail)
	err := s.cache.Get(ctx, key, nil)
	if errors.Is(err, cache.ErrCacheMiss) {
		return echo.NewHTTPError(
			http.StatusForbidden,
			"email not whitelisted",
		)
	}

	hashed_password, err := bcrypt.GenerateFromPassword(
		[]byte(req.NewPassword),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return echo.NewHTTPError(
			http.StatusInternalServerError,
			"invalid password",
		)
	}

	err = s.repo.UpdatePassword(
		ctx,
		req.InstituteEmail,
		string(hashed_password),
	)
	if err != nil {
		return echo.NewHTTPError(
			http.StatusInternalServerError,
			"failed to update password",
		)
	}

	ok, err := s.cache.Delete(ctx, key)
	if !ok || err != nil {
		return echo.NewHTTPError(
			http.StatusInternalServerError,
			"failed to delete email from whitelist",
		)
	}

	return nil
}
