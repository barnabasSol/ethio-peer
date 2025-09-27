package jwt

import "github.com/golang-jwt/jwt/v5"

type Claims struct {
	jwt.RegisteredClaims
	Roles    []string `json:"roles"`
	Username string   `json:"username"`
}
