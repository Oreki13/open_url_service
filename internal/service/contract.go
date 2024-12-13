package service

import (
	"context"
	"open_url_service/internal/appctx"
	"open_url_service/internal/entity"
)

type UserService interface {
	ListUser(ctx context.Context) (*[]entity.User, error)
	StoreUser(ctx context.Context) appctx.Response
}

type UrlService interface {
	FindUrlByPath(ctx context.Context, path string) (*entity.Url, error)
}
