package login

import (
	"context"
	"encoding/json"
	broker "ep-auth-service/internal/broker/rabbitmq"
	"ep-auth-service/internal/features/common"
	"ep-auth-service/internal/features/jwt"
	"ep-auth-service/internal/features/otp"
	"errors"
	"log"

	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	LoginUser(ctx context.Context, req LoginRequest) (*common.Response[LoginResponse], error)
}

type service struct {
	token_service jwt.Generator
	broker        *broker.RabbitMQ
	rep           Repository
	otp_manager   *otp.OTPManager
}

func NewService(
	rep Repository,
	rmq *broker.RabbitMQ,
	tg jwt.Generator,
	otp_manager *otp.OTPManager,
) Service {
	return &service{
		rep:           rep,
		token_service: tg,
		broker:        rmq,
		otp_manager:   otp_manager,
	}
}

func (s *service) LoginUser(
	ctx context.Context,
	login LoginRequest,
) (*common.Response[LoginResponse], error) {
	user, err := s.rep.GetUser(ctx, login)
	if err != nil {
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(login.Password),
	)

	if err != nil {
		return nil, ErrIncorrectPassword
	}

	if !user.InstituteEmailVerified {
		new_otp, err := s.otp_manager.Generate(user.Id.Hex())
		if err != nil {
			return nil, err
		}
		otp := broker.OtpPayload{
			Email: user.InstituteEmail,
			OTP:   new_otp.Value,
		}

		otp_json, err := json.Marshal(otp)
		if err != nil {
			return nil, errors.New("failed to marshal otp")
		}

		s.broker.Publish(broker.Message{
			Exchange: "notification_exchange",
			Topic:    "email.otp",
			Data:     otp_json,
		})

		return &common.Response[LoginResponse]{
			Message: "please verify your email",
			Data: LoginResponse{
				VerificationRequired: true,
				OtpSessionId:         &new_otp.SessionId,
			},
		}, nil
	}

	token, err := s.token_service.GenerateAccessToken(*user)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	refresh, err := s.token_service.GenerateRefreshToken(32)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	if err := s.rep.InsertRefreshToken(ctx, user.Id, refresh); err != nil {
		return nil, err
	}
	id := user.Id.Hex()
	return &common.Response[LoginResponse]{
		Message: "successfully logged in",
		Data: LoginResponse{
			VerificationRequired: false,
			UserId:               &id,
			AccessToken:          &token,
			RefreshToken:         &refresh,
		},
	}, nil

}
