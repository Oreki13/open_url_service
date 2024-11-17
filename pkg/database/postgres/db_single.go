package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	// check in runtime implement Databaser
	_ Adapter = (*DB)(nil)
)

type DB struct {
	//instanceID string
	txPool *pgxpool.Tx
	tx     *pgx.Tx
	db     *pgxpool.Pool // the Conn of the Tx, when tx != nil
	//opts       sql.TxOptions // valid when tx != nil
	reaMode bool
	dbName  string
}

func New(pool *pgxpool.Pool, readMode bool, dbName string) *DB {
	return &DB{
		db:      pool,
		reaMode: readMode,
		dbName:  dbName,
	}
}

func NewTx(pool *pgxpool.Pool, tx *pgx.Tx, readMode bool, dbName string) *DB {
	return &DB{
		db:      pool,
		tx:      tx,
		reaMode: readMode,
		dbName:  dbName,
	}
}

func (db *DB) Ping() error {
	return db.db.Ping(context.Background())
}

func (db *DB) InTransaction() bool {
	return db.tx != nil
}

// Close closes the database connection.
func (db *DB) Close() {
	db.db.Close()
}

// Exec executes a SQL statement and returns the number of rows it affected.
func (db *DB) Exec(ctx context.Context, query string, args ...any) (_ int64, err error) {
	if db.reaMode {
		return 0, fmt.Errorf("database mode read only")
	}

	res, err := db.execResult(ctx, query, args...)
	if err != nil {
		return 0, err
	}
	n := res.RowsAffected()

	return n, nil
}

func (db *DB) execResult(ctx context.Context, query string, args ...any) (formats pgconn.CommandTag, err error) {
	if db.tx != nil {
		return db.txPool.Exec(ctx, query, args...)
	}

	return db.db.Exec(ctx, query, args...)
}

// Query runs the DB query.
func (db *DB) Query(ctx context.Context, query string, args ...any) (pgx.Rows, error) {
	if db.tx != nil {
		return db.txPool.Query(ctx, query, args...)
	}

	return db.db.Query(ctx, query, args...)
}

// QueryRow runs the query and returns a single row.
func (db *DB) QueryRow(ctx context.Context, query string, args ...any) (pgx.Row, error) {

	if db.tx != nil {
		return db.txPool.QueryRow(ctx, query, args...), nil
	}

	return db.db.QueryRow(ctx, query, args...), nil
}

//// QueryX runs the DB query.
//func (db *DB) QueryX(ctx context.Context, query string, args ...any) (pgx.Rows, error) {
//	if db.tx != nil {
//		return db.tx.QueryContext(ctx, query, args...)
//	}
//
//	return db.db.QueryContext(ctx, query, args...)
//}
//
//// QueryRowX runs the query and returns a single row.
//func (db *DB) QueryRowX(ctx context.Context, query string, args ...any) pgx.Row {
//	if db.tx != nil {
//		return db.tx.QueryRowContext(ctx, query, args...)
//	}
//
//	return db.db.QueryRowContext(ctx, query, args...)
//}

// Transact executes the given function in the context of a SQL transaction at
// the given isolation level
func (db *DB) Transact(ctx context.Context, txFunc func(*DB) error) (err error) {
	if db.reaMode {
		return fmt.Errorf("database mode read only")
	}

	// For the levels which require retry, see
	// https://www.postgresql.org/docs/11/transaction-iso.html.

	return db.transact(ctx, txFunc)
}

func (db *DB) transact(ctx context.Context, txFunc func(*DB) error) (err error) {
	if db.InTransaction() {
		return errors.New("db transact function was called on a DB already in a transaction")
	}

	tx, err := db.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("tx begin: %w", err)
	}

	//defer func() {
	//	if p := recover(); p != nil {
	//		tx.Rollback()
	//	} else if err != nil {
	//		tx.Rollback()
	//	} else {
	//		if txErr := tx.Commit(); txErr != nil {
	//			err = fmt.Errorf("tx commit: %w", txErr)
	//		}
	//	}
	//}()

	dbtx := NewTx(db.db, &tx, false, db.dbName)

	if err := txFunc(dbtx); err != nil {
		err := tx.Rollback(ctx)
		if err != nil {
			return err
		}
		return fmt.Errorf("fn(tx): %w", err)
	}

	return tx.Commit(ctx)
}

// BeginTx start new transaction session
func (db *DB) BeginTx(ctx context.Context, opts pgx.TxOptions) (pgx.Tx, error) {
	return db.db.BeginTx(ctx, opts)
}
