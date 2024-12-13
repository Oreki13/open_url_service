package url

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
	"open_url_service/internal/appctx"
	"open_url_service/internal/controller/contract"
	"open_url_service/internal/service"
	"open_url_service/pkg/helper"
	"strings"
)

type findEndpointByPath struct {
	service service.UrlService
}

func (u *findEndpointByPath) Serve(xCtx appctx.Data) appctx.Response {
	ctx, span := helper.InitialOtelSpan(xCtx)
	defer span.End()

	title := xCtx.FiberCtx.Params("title")
	if strings.Contains(title, ".") {
		return *appctx.NewResponse().WithCode(404).WithMessage("Not Found!")
	}
	if len(title) == 0 {
		return *appctx.NewResponse().WithCode(404).WithMessage("Not Found!").WithData("No Title")
	}

	listUrl, err := u.service.FindUrlByPath(ctx, title)
	if err != nil {
		if errors.Is(pgx.ErrNoRows, err) {
			return *appctx.NewResponse().WithCode(404).WithMessage("Not Found!")
		}
		return *appctx.NewResponse().WithError([]appctx.ErrorResp{
			{
				Key:      "SERVER_ERROR",
				Messages: []string{err.Error()},
			},
		}).WithCode(fiber.StatusBadRequest)
	}

	return *appctx.NewResponse().WithCode(301).WithState(listUrl.Destination)
}

func NewFindEndpointByPath(svc service.UrlService) contract.Controller {
	return &findEndpointByPath{service: svc}
}
