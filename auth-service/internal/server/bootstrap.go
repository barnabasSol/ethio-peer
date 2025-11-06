package server

import (
	"context"
	"ep-auth-service/internal/features/common/cache"
	"ep-auth-service/internal/features/common/jwt"
	"ep-auth-service/internal/features/common/otp"
	"ep-auth-service/internal/features/login"
	refreshtoken "ep-auth-service/internal/features/refresh-token"
	resetpassword "ep-auth-service/internal/features/reset-password"
	"ep-auth-service/internal/features/signup"
	"ep-auth-service/internal/features/user"
	"log"
	"os"
	"time"
)

func (s *Server) bootstrap() error {
	redis_addr := os.Getenv("REDIS_ADDR")
	redis, err := cache.New(redis_addr, time.Minute*10)
	if err != nil {
		log.Fatal("failed to initialize redis cache")
		return err
	}
	log.Println("connected to redis")
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
		redis,
		token_gen,
	)
	otp.InitHandler(otp_service, s.echo.Group("/otp"))

	ref_repo := refreshtoken.NewRepository(s.db)
	ref_service := refreshtoken.NewService(
		ref_repo,
		s.broker,
		token_gen,
	)
	refreshtoken.InitHandler(ref_service, s.echo.Group("/refresh"))

	ur := user.NewRepository(s.db)
	us := user.NewService(ur)
	user.InitHandler(s.echo.Group("/admin"), us)

	rpr := resetpassword.NewRepository(s.db)
	rps := resetpassword.NewService(
		rpr,
		redis,
		s.broker,
		ur,
		otp_manager,
	)

	resetpassword.InitHandler(s.echo, rps)
	go func() {
		if err := s.g.Run(us); err != nil {
			log.Fatalf("%s", err.Error())
		}
	}()

	return nil
}
