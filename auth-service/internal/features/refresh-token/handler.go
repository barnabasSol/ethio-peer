package refreshtoken

import (
	"crypto/rsa"
	"ep-auth-service/internal/features/common"
	"errors"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	group *echo.Group
	s     Service
}

func InitHandler(s Service, group *echo.Group) *Handler {
	h := &Handler{
		group: group,
		s:     s,
	}
	h.group.GET("", h.Refresh)
	return h
}
func (h *Handler) Refresh(ctx echo.Context) error {
	withCookie := ctx.QueryParam("with_cookie")

	refreshTokenCookie, err := ctx.Cookie("refresh_token")
	if err != nil || refreshTokenCookie == nil || refreshTokenCookie.Value == "" {
		return echo.NewHTTPError(
			http.StatusUnauthorized,
			"failed to refresh, no refresh token provided",
		)
	}
	refreshToken := refreshTokenCookie.Value

	accessTokenCookie, err := ctx.Cookie("access_token")
	if err != nil || accessTokenCookie == nil || accessTokenCookie.Value == "" {
		return echo.NewHTTPError(
			http.StatusUnauthorized,
			"failed to refresh, no access token provided",
		)
	}
	accessToken := accessTokenCookie.Value

	publicKey, err := loadPublicKey("/root/certs/public.pem")
	if err != nil {
		ctx.Logger().Errorf("Failed to load public key: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "server error")
	}

	token, err := jwt.ParseWithClaims(
		accessToken,
		jwt.MapClaims{},
		func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, echo.NewHTTPError(
					http.StatusUnauthorized,
					"unexpected signing method",
				)
			}
			return publicKey, nil
		},
		jwt.WithValidMethods([]string{"RS256"}),
	)

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {

		} else {
			ctx.Logger().Errorf("Failed to validate access token: %v", err)
			return echo.NewHTTPError(
				http.StatusUnauthorized,
				"invalid or tampered access token",
			)
		}
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid access token claims")
	}

	userID, ok := claims["sub"].(string)
	if !ok || userID == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, "missing or invalid sub claim")
	}

	ctx.Logger().Infof("UserID: %s", userID)

	req := Request{
		UserId:       userID,
		RefreshToken: refreshToken,
	}

	result, err := h.s.Refresh(ctx.Request().Context(), req)
	if err != nil {
		return err
	}

	if withCookie == "true" && result.Data.AccessToken != "" && result.Data.RefreshToken != "" {
		atc := common.SetCookie(
			"access_token",
			result.Data.AccessToken,
			15,
		)
		rtc := common.SetCookie(
			"refresh_token",
			result.Data.RefreshToken,
			60*24*7,
		)
		ctx.SetCookie(atc)
		ctx.SetCookie(rtc)
		return ctx.JSON(http.StatusOK, common.Response[struct{}]{
			Message: result.Message,
			Data:    struct{}{},
		})
	}

	return ctx.JSON(http.StatusOK, result)
}

func loadPublicKey(path string) (*rsa.PublicKey, error) {
	pemData, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	key, err := jwt.ParseRSAPublicKeyFromPEM(pemData)
	if err != nil {
		return nil, err
	}
	return key, nil
}
