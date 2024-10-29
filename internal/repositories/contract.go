package repositories

import (
	"context"
	"database/sql"
	"open_url_service/internal/entity"
)

type UserRepository interface {
	ListUser(ctx context.Context) (*[]entity.User, error)
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
	Store(ctx context.Context, payload entity.User, opts ...Option) (int, error)
}
