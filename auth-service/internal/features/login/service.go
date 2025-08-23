package login

import "ep-auth-service/internal/features/jwt"

type Service struct {
	token_service jwt.Generator
}
