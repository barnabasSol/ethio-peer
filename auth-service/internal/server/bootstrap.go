package server

import (
	"ep-auth-service/internal/features/jwt"
	"ep-auth-service/internal/features/signup"
	"log"
)

func (s *Server) bootstrap() error {

	auth_group := s.echo.Group("auth")

	token_gen, err := jwt.NewTokenGenerator()

	if err != nil {
		log.Println("failed to initialize token generator")
		return err
	}

	signup_repo := signup.NewRepository(s.db)

	signup_service := signup.NewService(signup_repo, token_gen)

	signup.InitHandler(signup_service, auth_group)

	return nil
}
