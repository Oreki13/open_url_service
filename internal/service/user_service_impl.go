package service

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"open_url_service/internal/appctx"
	"open_url_service/internal/entity"
	"open_url_service/internal/repositories"
	"open_url_service/pkg/logger"
	"open_url_service/pkg/tracer"
	"time"
)

type userServiceImpl struct {
	repo repositories.UserRepository
}

func (u userServiceImpl) ListUser(ctx context.Context) (*[]entity.User, error) {
	return u.repo.ListUser(ctx)
}

func (u userServiceImpl) StoreUser(ctx context.Context) appctx.Response {
	var (
		lf = logger.NewFields(
			logger.EventName("ServiceStoreUser"),
		)
	)
	ctx, span := tracer.NewSpan(ctx, "Service.StoreUser", nil)
	defer span.End()

	payload := entity.User{
		ID:   1,
		Name: "John Doe",
		Age:  12,
	}

	// value custom logger
	lf.Append(logger.Any("payload.id", payload.ID))
	lf.Append(logger.Any("payload.name", payload.Name))
	lf.Append(logger.Any("payload.age", payload.Age))

	// start db transaction
	tx, err := u.repo.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelSerializable,
	})

	if err != nil {
		tracer.AddSpanError(span, err)
		logger.ErrorWithContext(ctx, fmt.Sprintf("start db transaction got error: %v", err), lf...)
		return *appctx.NewResponse().WithCode(http.StatusInternalServerError).WithMessage("Somethin went wrong")
	}

	txOpt := repositories.WithTransaction(tx)

	_, err = u.repo.Store(ctx, entity.User{
		Name: "John Doe",
		Age:  12,
	}, txOpt)
	if err != nil {
		tracer.AddSpanError(span, err)
		logger.ErrorWithContext(ctx, fmt.Sprintf("store user got error: %v", err), lf...)

		// rollback transaction
		err := tx.Rollback()
		if err != nil {
			return appctx.Response{}
		}
		return *appctx.NewResponse().WithCode(http.StatusInternalServerError).WithMessage("Somethin went wrong")
	}

	logger.InfoWithContext(ctx, "success store user", lf...)

	// commit transaction
	err = tx.Commit()
	if err != nil {
		return appctx.Response{}
	}
	return *appctx.NewResponse().WithCode(http.StatusCreated).WithMessage("Success created user").WithData(
		map[string]interface{}{
			"user_name":  payload.Name,
			"created_at": time.Now(),
		},
	)
}

func NewUserServiceImpl(repo repositories.UserRepository) UserService {
	return &userServiceImpl{
		repo: repo,
	}
}
