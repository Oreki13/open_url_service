package repositories

import (
	"github.com/jackc/pgx/v5"
)

type option struct {
	tx pgx.Tx
}

type Option func(*option)

func WithTransaction(tx pgx.Tx) Option {
	return func(o *option) {
		o.tx = tx
	}
}
