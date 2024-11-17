package service

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
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
		ID:        "id",
		Name:      "John Doe",
		Email:     "john@mail.com",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// value custom logger
	lf.Append(logger.Any("payload.id", payload.ID))
	lf.Append(logger.Any("payload.name", payload.Name))
	lf.Append(logger.Any("payload.age", payload.Email))

	// start db transaction
	tx, err := u.repo.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.Serializable,
	})

	if err != nil {
		tracer.AddSpanError(span, err)
		logger.ErrorWithContext(ctx, fmt.Sprintf("start db transaction got error: %v", err), lf...)
		return *appctx.NewResponse().WithCode(http.StatusInternalServerError).WithMessage("Somethin went wrong")
	}

	txOpt := repositories.WithTransaction(tx)

	_, err = u.repo.Store(ctx, entity.User{
		Name:      "John Doe",
		Email:     "john@mail.com",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, txOpt)
	if err != nil {
		tracer.AddSpanError(span, err)
		logger.ErrorWithContext(ctx, fmt.Sprintf("store user got error: %v", err), lf...)

		// rollback transaction
		err := tx.Rollback(ctx)
		if err != nil {
			return appctx.Response{}
		}
		return *appctx.NewResponse().WithCode(http.StatusInternalServerError).WithMessage("Somethin went wrong")
	}

	logger.InfoWithContext(ctx, "success store user", lf...)

	// commit transaction
	err = tx.Commit(ctx)
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
