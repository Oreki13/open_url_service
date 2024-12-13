package service

import (
	"context"
	"open_url_service/internal/entity"
	"open_url_service/internal/repositories"
	"open_url_service/pkg/tracer"
)

type urlServiceImpl struct {
	repo repositories.UrlRepository
}

func (u *urlServiceImpl) FindUrlByPath(appCtx context.Context, path string) (*entity.Url, error) {
	ctx, span := tracer.NewSpan(appCtx, "Service.FindUrlByPath", nil)
	defer span.End()
	return u.repo.FindUrlByPath(ctx, path)
}

func NewFindUrlByPath(repo repositories.UrlRepository) UrlService {
	return &urlServiceImpl{
		repo: repo,
	}
}
