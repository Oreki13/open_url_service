package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"open_url_service/pkg/config"
)

func connect(cnf *config.Config) (*pgxpool.Pool, error) {
	var (
		err error
	)

	conf, err := NewPgsqlConfig(cnf)

	if err != nil {
		return nil, err
	}

	db, err := pgxpool.NewWithConfig(context.Background(), conf)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(context.Background()); err != nil {
		db.Close()
		return nil, err
	}

	//db.SetMaxIdleConns(dbConfig.MaxIdleConn)
	//db.SetMaxOpenConns(dbConfig.MaxOpenConn)
	//db.SetConnMaxLifetime(time.Duration(dbConfig.MaxConnLifetime) * time.Hour)
	//db.SetConnMaxIdleTime(time.Duration(dbConfig.MaxIdleTime) * time.Hour)

	return db, nil
}

func NewPgsqlConfig(cnf *config.Config) (*pgxpool.Config, error) {
	dbConfig := cnf.DatabaseConfig

	urlScheme := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?",
		dbConfig.DBUser, dbConfig.DBPassword, dbConfig.DBHost, dbConfig.DBPort, dbConfig.DBName)

	pgConfig, err := pgxpool.ParseConfig(urlScheme)

	if err != nil {
		return nil, err
	}

	//pgConfig.BeforeAcquire = func(ctx context.Context, c *pgx.Conn) bool {
	//	return true
	//}
	//
	//pgConfig.AfterRelease = func(c *pgx.Conn) bool {
	//	logger.Info("After releasing the connection pool to the database!!")
	//	return true
	//}
	//
	//pgConfig.BeforeClose = func(c *pgx.Conn) {
	//	logger.Info("Closed the connection pool to the database!!")
	//}

	//if dbConfig.TLS {
	//	urlScheme += fmt.Sprintf("sslmode=require sslrootcert=%s sslcert=%s sslkey=%s",
	//		dbConfig.CAPath, dbConfig.ClientCertPath, dbConfig.ClientKeyPath)
	//}
	//tlsConfig, err := dbConfig.TlsConfig(cnf.AppEnv)
	//if err != nil {
	//	return nil, err
	//}

	//if tlsConfig != nil {
	//	if err = postgres.RegisterTLSConfig("custom", tlsConfig); err != nil {
	//		return nil, err
	//	}
	//
	//	//conf.TLSConfig = "custom"
	//}

	return pgConfig, nil
}

func ConnectDatabase(cnf *config.Config) (*pgxpool.Pool, error) {
	db, err := connect(cnf)

	if err != nil {
		return nil, err
	}
	return db, nil
}
