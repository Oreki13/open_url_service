package user

import (
	"open_url_service/internal/appctx"
	"open_url_service/internal/controller/contract"
	"open_url_service/internal/service"
	"open_url_service/pkg/tracer"
)

type storeUser struct {
	service service.UserService
}

func (g *storeUser) Serve(xCtx appctx.Data) appctx.Response {
	ctx, span := tracer.NewSpan(xCtx.FiberCtx.Context(), "Controller.CreateUser", nil)
	defer span.End()

	res := g.service.StoreUser(ctx)
	return res
}

func NewStoreUser(svc service.UserService) contract.Controller {
	return &storeUser{service: svc}
}
