package login

import (
	"context"
	broker "ep-auth-service/internal/broker/rabbitmq"
	"ep-auth-service/internal/features/jwt"
	"ep-auth-service/internal/features/otp"
	"ep-auth-service/internal/features/shared"
	"log"

	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	LoginUser(ctx context.Context, req LoginRequest) (*shared.Response[LoginReponse], error)
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
) (*shared.Response[LoginReponse], error) {
	user, err := s.rep.GetUser(ctx, login)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(login.Password))

	if err != nil {
		return nil, ErrIncorrectPassword
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

	return &shared.Response[LoginReponse]{
		Message: "successfully logged in",
		Data: LoginReponse{
			UserId:       user.Id.Hex(),
			AccessToken:  token,
			RefreshToken: refresh,
		},
	}, nil

}
