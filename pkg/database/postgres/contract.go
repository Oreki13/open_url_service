package postgres

import (
	"context"
	"github.com/jackc/pgx/v5"
)

type Adapter interface {
	Ping() error
	InTransaction() bool
	Close()
	Query(ctx context.Context, query string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, query string, args ...interface{}) (pgx.Row, error)
	//QueryX(ctx context.Context, query string, args ...any) (pgx.Rows, error)
	//QueryRowX(ctx context.Context, query string, args ...any) pgx.Row
	Exec(ctx context.Context, query string, args ...interface{}) (_ int64, err error)
	Transact(ctx context.Context, txFunc func(*DB) error) (err error)
	BeginTx(ctx context.Context, opts pgx.TxOptions) (pgx.Tx, error)
}
