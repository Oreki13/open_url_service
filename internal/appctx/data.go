package appctx

import (
	"github.com/gofiber/fiber/v2"
	"open_url_service/pkg/config"
)

type Data struct {
	FiberCtx *fiber.Ctx
	Cfg      *config.Config
}
