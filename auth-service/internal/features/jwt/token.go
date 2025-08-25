package jwt

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"ep-auth-service/internal/models"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Generator interface {
	GenerateAccessToken(user models.User) (string, error)
	GenerateRefreshToken(n int) (string, error)
}

type generator struct {
	privateKey *rsa.PrivateKey
}

func NewTokenGenerator() (Generator, error) {
	pk_pem, err := ReadPrivateKey()
	if err != nil {
		return nil, err
	}

	privKey, err := jwt.ParseRSAPrivateKeyFromPEM(pk_pem)

	if err != nil {
		return nil, err
	}

	return generator{
		privateKey: privKey,
	}, nil
}

func (g generator) GenerateAccessToken(user models.User) (string, error) {

	expiry, err := strconv.Atoi(os.Getenv("JWT_EXPIRY_MINS"))
	if err != nil {
		return "", err
	}

	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.NewString(),
			Subject:   user.Id.Hex(),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expiry) * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "ep.barney-host.site",
			Audience:  []string{"ep-web.barney-host.site"},
		},
		Roles:    user.Roles,
		Username: user.Username,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	token.Header["kid"] = "7a101016-033e-44de-9137-572a113a592f"

	signed, err := token.SignedString(g.privateKey)
	if err != nil {
		return "", err
	}

	return signed, nil
}

func (g generator) GenerateRefreshToken(n int) (string, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
