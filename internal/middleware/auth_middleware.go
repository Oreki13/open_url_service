package middleware

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"open_url_service/internal/appctx"
	"open_url_service/pkg/config"
	"strings"
)

type AuthMiddleware struct {
}

func NewAuthMiddleware() *AuthMiddleware {
	return &AuthMiddleware{}
}

func (a *AuthMiddleware) Authenticate(xCtx *fiber.Ctx, conf *config.Config) appctx.Response {
	auth := xCtx.GetReqHeaders()["Authorization"]
	userId := xCtx.GetReqHeaders()["X-Control-User"]
	secretKey := []byte(conf.SecretKey)

	if len(auth) == 0 || len(userId) == 0 {
		return *appctx.NewResponse().WithCode(fiber.StatusUnauthorized).WithMessage("Unauthorized no data")
	}

	trimAuth := strings.TrimPrefix(auth[0], "Bearer ")

	token, err := jwt.Parse(trimAuth, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		switch {
		case errors.Is(err, jwt.ErrTokenExpired):
			return *appctx.NewResponse().WithCode(fiber.StatusUnauthorized).WithMessage("Token has been expired")
		case errors.Is(err, jwt.ErrTokenMalformed):
			return *appctx.NewResponse().WithCode(fiber.StatusUnauthorized).WithMessage("Token is invalid")
		}

		return *appctx.NewResponse().WithCode(fiber.StatusUnauthorized).WithMessage("Unauthorized")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if claims["id"] != userId[0] {
			return *appctx.NewResponse().WithCode(fiber.StatusUnauthorized).WithMessage("Invalid token user")

		}
		return *appctx.NewResponse().WithCode(fiber.StatusOK)
	}
	return *appctx.NewResponse().WithCode(fiber.StatusUnauthorized).WithMessage("Unauthorized claim")

}
