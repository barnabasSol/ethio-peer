package jwt

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"ep-auth-service/internal/models"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Generator interface {
	ParseRefreshTokenJWT(tokenString string) (string, error)
	GenerateAccessToken(user models.User) (string, error)
	GenerateRefreshToken(n int) (string, error)
	GenerateRefreshTokenJWT(user_id string) (string, error)
}

type generator struct {
	privateKey *rsa.PrivateKey
}

func NewTokenGenerator() (Generator, error) {
	pk_pem, err := ReadPrivateKey()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	privKey, err := jwt.ParseRSAPrivateKeyFromPEM(pk_pem)

	if err != nil {
		log.Println(err)
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
			Audience: []string{
				"https://ep-web.barney-host.site",
				"http://localhost:5173",
			},
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

func (g generator) GenerateRefreshTokenJWT(userId string) (string, error) {
	rt_exp := os.Getenv("RT_EXPIRY_DAYS")
	if rt_exp == "" {
		rt_exp = "7"
	}
	exp, err := strconv.Atoi(rt_exp)
	if err != nil {

	}
	claims := jwt.RegisteredClaims{
		Subject:   userId,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(exp) * time.Minute)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token.Header["kid"] = "7a101016-033e-44de-9137-572a113a592f"

	signed, err := token.SignedString(g.privateKey)
	if err != nil {
		log.Println(err)
		return "", errors.New("failed to sign token")
	}
	return signed, nil
}

func (g generator) ParseRefreshTokenJWT(tokenString string) (string, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&jwt.RegisteredClaims{},
		func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fmt.Errorf(
					"unexpected signing method: %v",
					token.Header["alg"],
				)
			}
			return &g.privateKey.PublicKey, nil
		},
	)

	if err != nil {
		return "", errors.New("tampered token")
	}

	if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok && token.Valid {
		return claims.Subject, nil
	}
	return "", errors.New("invalid token claims or token expired")
}
