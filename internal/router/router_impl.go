package router

import (
	"github.com/gofiber/fiber/v2"
	"open_url_service/internal/appctx"
	"open_url_service/internal/bootstrap"
	"open_url_service/internal/controller"
	"open_url_service/internal/controller/contract"
	"open_url_service/internal/controller/url"
	"open_url_service/internal/controller/user"
	"open_url_service/internal/handler"
	"open_url_service/internal/middleware"
	"open_url_service/internal/repositories"
	"open_url_service/internal/service"
	"open_url_service/pkg/config"
)

type router struct {
	cfg   *config.Config
	fiber fiber.Router
}

func (rtr *router) handle(hfn httpHandlerFunc, svc contract.Controller, mdws ...middleware.MiddlewareFunc) fiber.Handler {
	return func(xCtx *fiber.Ctx) error {
		//check registered middleware functions
		if rm := middleware.FilterFunc(rtr.cfg, xCtx, mdws); rm.Code != fiber.StatusOK {
			// return response base on middleware
			res := *appctx.NewResponse().
				WithCode(rm.Code).
				WithError(rm.Errors).
				WithMessage(rm.Message)
			return rtr.response(xCtx, res)
		}

		//send to controller
		resp := hfn(xCtx, svc, rtr.cfg)
		return rtr.response(xCtx, resp)
	}
}

func (rtr *router) response(fiberCtx *fiber.Ctx, resp appctx.Response) error {
	if resp.Code == 301 && len(resp.State) != 0 {
		return fiberCtx.Redirect(resp.State, 301)
	}
	fiberCtx.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	return fiberCtx.Status(resp.Code).Send(resp.Byte())
}

func (rtr *router) Route() {
	//init db
	db := bootstrap.RegistryDatabase(rtr.cfg)

	//define repositories
	userRepo := repositories.NewUserRepositoryImpl(db)
	urlRepo := repositories.NewFindUrlByPathImpl(db)

	//define services
	userSvc := service.NewUserServiceImpl(userRepo)
	urlSvc := service.NewFindUrlByPath(urlRepo)

	//define middleware
	basicMiddleware := middleware.NewAuthMiddleware()

	//define provider
	//example := provider.NewExampleProvider(rtr.cfg)

	//define storage
	//fs := bootstrap.RegistryGCS(rtr.cfg.GCSConfig.ServiceAccountPath)
	//fs := bootstrap.RegistryAWSSession(rtr.cfg)
	//fs := bootstrap.RegistryMinio(rtr.cfg)

	//define controller
	getAllUser := user.NewGetAllUser(userSvc)
	storeUser := user.NewStoreUser(userSvc)
	findUrlByPath := url.NewFindEndpointByPath(urlSvc)
	//getTodos := todo.NewGetTodo(example)

	health := controller.NewGetHealth()
	internalV1 := rtr.fiber.Group("/api/internal/v1")

	rtr.fiber.Get("/:title", rtr.handle(
		handler.HttpRequest, findUrlByPath))

	rtr.fiber.Get("/ping", rtr.handle(
		handler.HttpRequest,
		health,
	))

	internalV1.Get("/users", rtr.handle(
		handler.HttpRequest,
		getAllUser,
		//middleware
		basicMiddleware.Authenticate,
	))

	internalV1.Post("/users", rtr.handle(
		handler.HttpRequest,
		storeUser,
		//middleware
		// basicMiddleware.Authenticate,
	))

	//internalV1.Get("/todos", rtr.handle(
	//	handler.HttpRequest,
	//	getTodos,
	//	//middleware
	//	// basicMiddleware.Authenticate,
	//))

}

func NewRouter(cfg *config.Config, fiber fiber.Router) Router {
	return &router{
		cfg:   cfg,
		fiber: fiber,
	}
}
