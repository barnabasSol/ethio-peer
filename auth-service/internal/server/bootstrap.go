package server

import (
	"context"
	"ep-auth-service/internal/features/jwt"
	"ep-auth-service/internal/features/login"
	"ep-auth-service/internal/features/otp"
	"ep-auth-service/internal/features/signup"
	"log"
)

func (s *Server) bootstrap() error {
	auth_group := s.echo.Group("")
	token_gen, err := jwt.NewTokenGenerator()
	if err != nil {
		log.Println("failed to initialize token generator")
		return err
	}

	otp_manager := otp.NewOTPManager(context.Background())
	signup_repo := signup.NewRepository(s.db)
	signup_service := signup.NewService(
		signup_repo,
		s.broker,
		token_gen,
		otp_manager,
	)

	signup.InitHandler(signup_service, auth_group)
	login_repo := login.NewRepository(s.db)
	login_service := login.NewService(
		login_repo,
		s.broker,
		token_gen,
		otp_manager,
	)
	login.InitHandler(login_service, auth_group)

	otp_repo := otp.NewRepository(s.db)
	otp_service := otp.NewService(
		otp_manager,
		otp_repo,
		token_gen,
	)
	otp.InitHandler(otp_service, s.echo.Group("/otp"))
	return nil
}
