package middleware

import (
	"github.com/gofiber/fiber/v2"
	"open_url_service/internal/appctx"
	"open_url_service/pkg/config"
)

type MiddlewareFunc func(xCtx *fiber.Ctx, conf *config.Config) appctx.Response

// FilterFunc is a iterator resolver in each middleware registered
func FilterFunc(conf *config.Config, xCtx *fiber.Ctx, mfs []MiddlewareFunc) appctx.Response {
	// Initiate postive case
	var response = appctx.Response{Code: fiber.StatusOK}
	for _, mf := range mfs {
		if response = mf(xCtx, conf); response.Code != fiber.StatusOK {
			return response
		}
	}

	return response
}
