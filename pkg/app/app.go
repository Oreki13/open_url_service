package app

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/google/uuid"
	"open_url_service/internal/bootstrap"
	"open_url_service/internal/router"
	"open_url_service/pkg/config"
	"strings"
)

type App struct {
	*fiber.App
	Cfg *config.Config
}

var appServer *App

func InitializeApp(cfg *config.Config) {
	// boostrap run and initialize package dependency
	bootstrap.RegistryLogger(cfg)
	otelManager, err := bootstrap.RegistryOpenTelemetry(cfg)

	if err != nil {

		err := otelManager.Close()
		if err != nil {
			return
		}
	}

	f := fiber.New(cfg.FiberConfig())
	f.Use(
		cors.New(cors.Config{
			MaxAge: 300,
			AllowOrigins: strings.Join([]string{
				"*",
			}, ","),
			AllowHeaders: strings.Join([]string{
				"Origin",
				"Content-Type",
				"Accept",
			}, ","),
			AllowMethods: strings.Join([]string{
				fiber.MethodGet,
				fiber.MethodPost,
			}, ","),
		}),
		requestid.New(requestid.Config{
			ContextKey: "refid",
			Header:     "X-Reference-Id",
			Generator: func() string {
				return uuid.New().String()
			},
		}),
		logger.New(logger.Config{
			Format: "[${ip}]:${port} ${status} - ${method} ${path}\n",
		}),
	)

	rtr := router.NewRouter(cfg, f)
	rtr.Route()

	appServer = &App{
		App: f,
		Cfg: cfg,
	}
}

func (app *App) StartServer() (err error) {
	return app.Listen(fmt.Sprintf("%v:%v", app.Cfg.AppHost, app.Cfg.AppPort))
}

func GetServer() *App {
	return appServer
}
