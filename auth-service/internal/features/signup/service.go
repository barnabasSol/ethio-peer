package signup

import (
	"context"
	"ep-auth-service/internal/features/jwt"
	"ep-auth-service/internal/features/shared"
	"ep-auth-service/internal/models"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	SignUpUser(
		ctx context.Context,
		user SignUpRequest,
	) (
		*shared.Response[SignUpResponse],
		error,
	)
}

type service struct {
	repo          Repository
	token_service jwt.Generator
}

func NewService(
	repo Repository,
	ts jwt.Generator,
) Service {
	return &service{
		token_service: ts,
		repo:          repo,
	}
}

func (s *service) SignUpUser(
	ctx context.Context,
	user SignUpRequest,
) (
	*shared.Response[SignUpResponse],
	error,
) {

	hashed_password, err := bcrypt.GenerateFromPassword(
		[]byte(user.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return nil, err
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
		return nil, err
	}

	user_model.Id = id

	token, err := s.token_service.GenerateAccessToken(user_model)
	if err != nil {
		return nil, err
	}

	refresh, err := s.token_service.GenerateRefreshToken(32)
	if err != nil {
		return nil, err
	}

	return &shared.Response[SignUpResponse]{
		Message: "user successfully created",
		Data: SignUpResponse{
			UserId:       id.Hex(),
			AccessToken:  token,
			RefreshToken: refresh,
		},
	}, nil
}
