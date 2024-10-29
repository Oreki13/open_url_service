package middleware

import (
	"encoding/base64"
	"github.com/gofiber/fiber/v2"
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
	auth := xCtx.GetReqHeaders()["Authorization"][0]

	if len(auth) == 0 {
		return *appctx.NewResponse().WithCode(fiber.StatusUnauthorized).WithMessage("Unauthorized")
	}

	decodeString, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(auth, "Basic "))
	if err != nil {
		return *appctx.NewResponse().WithCode(fiber.StatusUnauthorized).WithMessage("Unauthorized")
	}

	resultAuth := strings.Split(string(decodeString), ":")

	if resultAuth[0] == "username" && resultAuth[1] == "password" {
		return *appctx.NewResponse().WithCode(fiber.StatusOK)
	}

	return *appctx.NewResponse().WithCode(fiber.StatusUnauthorized).WithMessage("Unauthorized")
}
