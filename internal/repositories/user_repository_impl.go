package repositories

import (
	"context"
	"github.com/jackc/pgx/v5"
	"open_url_service/internal/entity"
	"open_url_service/pkg/database/postgres"
	"open_url_service/pkg/helper"
	"open_url_service/pkg/logger"
	"open_url_service/pkg/tracer"
)

type userRepositoryImpl struct {
	db postgres.Adapter
}

func (u *userRepositoryImpl) ListUser(ctx context.Context) (*[]entity.User, error) {
	query := "SELECT u.id, u.name, u.email, r.name, u.created_at, u.updated_at as role FROM users u INNER JOIN role_user r ON u.role_id = r.id WHERE is_deleted = 0"
	var results []entity.User
	logger.Debug("DEbug test")
	logger.Info("Info test")
	logger.Warn("warn test")
	logger.Error("Error test")
	//logger.Fatal("fatar test")
	res, err := u.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	for res.Next() {
		var result entity.User
		err = res.Scan(&result.ID, &result.Email, &result.Name, &result.Role, &result.CreatedAt, &result.UpdatedAt)
		if err != nil {
			logger.Error(err)
			return nil, err
		}

		results = append(results, result)
	}
	//sum(rate(log_messages_total{app="loki",level=~"error|warn"}[1m])) by (level)
	//{app="loki"} | logfmt | level="warn" or level="error"
	if res.Err() != nil {
		return nil, res.Err()
	}

	return &results, nil
}

func (r userRepositoryImpl) BeginTx(ctx context.Context, opts pgx.TxOptions) (pgx.Tx, error) {
	return r.db.BeginTx(ctx, opts)
}

func (r userRepositoryImpl) Store(ctx context.Context, payload entity.User, opts ...Option) (int, error) {
	var (
		id  int
		err error
		tx  pgx.Tx
	)

	ctx, span := tracer.NewSpan(ctx, "Repo.StoreUser", nil)
	defer span.End()

	opt := &option{}
	for _, f := range opts {
		f(opt)
	}

	if opt.tx != nil {
		tx = opt.tx

	} else {
		tx, err = r.db.BeginTx(ctx, pgx.TxOptions{
			IsoLevel: "serializable",
		})

		if err != nil {
			tracer.AddSpanError(span, err)
			return 0, err
		}

		defer func(tx pgx.Tx) {
			err := tx.Commit(ctx)
			if err != nil {

			}
		}(tx)
	}

	query, val, err := helper.StructQueryInsert(payload, "users", "db", false)

	rows, err := tx.Query(
		ctx,
		query,
		val...,
	)
	if err != nil {
		tracer.AddSpanError(span, err)
		return 0, err
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&id); err != nil {
			tracer.AddSpanError(span, err)
			return 0, err
		}
	}

	return id, err
}

func NewUserRepositoryImpl(db postgres.Adapter) UserRepository {
	return &userRepositoryImpl{
		db: db,
	}
}
