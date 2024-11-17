package user

import (
	"github.com/gofiber/fiber/v2"
	"open_url_service/internal/appctx"
	"open_url_service/internal/controller/contract"
	"open_url_service/internal/service"
	"open_url_service/pkg/tracer"
)

type getAllUser struct {
	service service.UserService
}

func (g *getAllUser) Serve(xCtx appctx.Data) appctx.Response {
	//ctx := xCtx.FiberCtx.Context()
	ctx, span := tracer.NewSpan(xCtx.FiberCtx.Context(), "Controller.getAllUser", nil)
	defer span.End()
	users, err := g.service.ListUser(ctx)
	if err != nil {
		return *appctx.NewResponse().WithError([]appctx.ErrorResp{
			{
				Key:      "PROVIDER_ERR",
				Messages: []string{err.Error()},
			},
		}).WithMessage(err.Error()).WithCode(fiber.StatusBadRequest)
	}

	return *appctx.NewResponse().WithCode(fiber.StatusOK).WithMessage("Success").WithData(users)
}

func NewGetAllUser(svc service.UserService) contract.Controller {
	return &getAllUser{service: svc}
}
