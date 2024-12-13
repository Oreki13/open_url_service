package repositories

import (
	"context"
	"fmt"
	"go.opentelemetry.io/otel/attribute"
	"open_url_service/internal/entity"
	"open_url_service/pkg/database/postgres"
	"open_url_service/pkg/helper"
	"open_url_service/pkg/logger"
	"open_url_service/pkg/tracer"
	"runtime/debug"
	"time"
)

type urlRepositoryImpl struct {
	db postgres.Adapter
}

func (u *urlRepositoryImpl) FindUrlByPath(appCtx context.Context, path string) (*entity.Url, error) {

	ctx, span := tracer.NewSpan(appCtx, "Repository.FindUrlByPath", nil)

	defer span.End()

	query := fmt.Sprintf("SELECT id, title, path, destination, count_clicks, updated_at, created_at FROM tbl_data_url WHERE  tbl_data_url.path = '%s'", path)

	start := time.Now()
	res, err := u.db.QueryRow(ctx, query)
	duration := time.Since(start)

	span.SetAttributes(attribute.String("DB.Query", query), attribute.Int64("DB.Duration_ms", duration.Milliseconds()))
	if err != nil {
		helper.SetOtelError(span, err, string(debug.Stack()))
		return nil, err
	}

	query = fmt.Sprintf("UPDATE tbl_data_url SET count_clicks = count_clicks + 1 WHERE path = '%s'", path)
	start = time.Now()
	_, err = u.db.Exec(ctx, query)
	duration = time.Since(start)

	span.SetAttributes(attribute.String("DB.Query.add", query), attribute.Int64("DB.Duration.add_ms", duration.Milliseconds()))
	if err != nil {
		helper.SetOtelError(span, err, string(debug.Stack()))
		return nil, err
	}

	var result entity.Url
	err = res.Scan(&result.ID, &result.Title, &result.Path, &result.Destination, &result.CountClick, &result.UpdatedAt, &result.CreatedAt)

	if err != nil {
		logger.Error(err)
		helper.SetOtelError(span, err, string(debug.Stack()))

		return nil, err
	}

	return &result, nil
}

func NewFindUrlByPathImpl(db postgres.Adapter) UrlRepository {
	return &urlRepositoryImpl{
		db: db,
	}
}
