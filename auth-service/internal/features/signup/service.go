package signup

import (
	"context"
	"encoding/json"
	broker "ep-auth-service/internal/broker/rabbitmq"
	"ep-auth-service/internal/features/common"
	"ep-auth-service/internal/features/common/jwt"
	"ep-auth-service/internal/features/common/otp"
	"ep-auth-service/internal/models"
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	SignUpUser(
		ctx context.Context,
		user SignUpRequest,
	) (
		*common.Response[SignUpResponse],
		error,
	)
}

type service struct {
	otp_manager   *otp.OTPManager
	repo          Repository
	broker        *broker.RabbitMQ
	token_service jwt.Generator
}

func NewService(
	repo Repository,
	broker *broker.RabbitMQ,
	ts jwt.Generator,
	otp_manager *otp.OTPManager,
) Service {
	return &service{
		token_service: ts,
		broker:        broker,
		repo:          repo,
		otp_manager:   otp_manager,
	}
}

func (s *service) SignUpUser(
	ctx context.Context,
	user SignUpRequest,
) (
	*common.Response[SignUpResponse],
	error,
) {
	hashed_password, err := bcrypt.GenerateFromPassword(
		[]byte(user.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		log.Println(err)
		return nil, echo.NewHTTPError(
			http.StatusInternalServerError,
			"failed to process password",
		)
	}
	user_model := models.User{
		Username:               user.Username,
		Name:                   user.Name,
		Roles:                  []string{"peer"},
		IsActive:               false,
		InstituteEmail:         user.InstituteEmail,
		PasswordHash:           string(hashed_password),
		InstituteEmailVerified: false,
		Email:                  user.Email,
		CreatedAt:              time.Now().UTC(),
		UpdatedAt:              time.Now().UTC(),
	}

	id, err := s.repo.Insert(ctx, user_model)

	if err != nil {
		log.Println(err)
		return nil, echo.NewHTTPError(
			http.StatusInternalServerError,
			"failed to create user",
		)
	}
	user_model.Id = id

	new_otp, err := s.otp_manager.Generate(id.Hex())

	if err != nil {
		return nil, echo.NewHTTPError(
			http.StatusInternalServerError,
			"failed to generate otp",
		)
	}

	otp := broker.OtpPayload{
		Email: user.InstituteEmail,
		OTP:   new_otp.Value,
	}

	new_peer := broker.PeerPayload{
		UserId:    id.Hex(),
		Interests: *user.Interests,
		Bio:       *user.Bio,
	}

	otp_json, err := json.Marshal(otp)
	if err != nil {
		return nil, echo.NewHTTPError(
			http.StatusInternalServerError,
			"failed to marshal otp",
		)
	}

	new_peer_json, err := json.Marshal(new_peer)
	if err != nil {
		return nil, echo.NewHTTPError(
			http.StatusInternalServerError,
			"something went wrong",
		)
	}

	s.broker.Publish(broker.Message{
		Exchange: "notification_exchange",
		Topic:    "email.otp",
		Data:     otp_json,
	})

	s.broker.Publish(broker.Message{
		Exchange: "new_peer_exchange",
		Topic:    "peer.new",
		Data:     new_peer_json,
	})

	return &common.Response[SignUpResponse]{
		Message: "please verify your email",
		Data: SignUpResponse{
			VerificationRequired: true,
			OtpSessionId:         &new_otp.SessionId,
		},
	}, nil

}
