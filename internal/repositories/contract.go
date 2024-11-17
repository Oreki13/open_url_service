package repositories

import (
	"context"
	"github.com/jackc/pgx/v5"
	"open_url_service/internal/entity"
)

type UserRepository interface {
	ListUser(ctx context.Context) (*[]entity.User, error)
	BeginTx(ctx context.Context, opts pgx.TxOptions) (pgx.Tx, error)
	Store(ctx context.Context, payload entity.User, opts ...Option) (int, error)
}
