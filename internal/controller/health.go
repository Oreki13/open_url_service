package controller

import (
	"github.com/gofiber/fiber/v2"
	"open_url_service/internal/appctx"
	"open_url_service/internal/controller/contract"
)

type getHealth struct {
}

func (g *getHealth) Serve(xCtx appctx.Data) appctx.Response {
	// Ping Endpoint
	return *appctx.NewResponse().WithCode(fiber.StatusOK).WithMessage("ok").WithData(struct {
		Message string `json:"message"`
	}{
		Message: "Waras!",
	})
}

func NewGetHealth() contract.Controller {
	return &getHealth{}
}
