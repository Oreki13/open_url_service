package mysql

import (
	"fmt"
	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"open_url_service/pkg/config"
	"open_url_service/pkg/logger"
	"time"
)

func connect(cnf *config.Config) (*sqlx.DB, error) {
	var (
		err      error
		dbConfig = cnf.DatabaseConfig
	)

	conf, err := NewMysqlConfig(cnf)
	if err != nil {
		logger.Fatal("Failed to create mysql config")
	}

	dsn := conf.FormatDSN()
	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		err := db.Close()
		if err != nil {
			return nil, err
		}
		return nil, err
	}

	db.SetMaxIdleConns(dbConfig.MaxIdleConn)
	db.SetMaxOpenConns(dbConfig.MaxOpenConn)
	db.SetConnMaxLifetime(time.Duration(dbConfig.MaxConnLifetime) * time.Hour)
	db.SetConnMaxIdleTime(time.Duration(dbConfig.MaxIdleTime) * time.Hour)

	return db, nil
}

func NewMysqlConfig(cnf *config.Config) (*mysql.Config, error) {
	dbConfig := cnf.DatabaseConfig
	conf := mysql.NewConfig()
	conf.Net = "tcp"
	conf.Addr = fmt.Sprintf("%v:%v", dbConfig.DBHost, dbConfig.DBPort)
	conf.User = dbConfig.DBUser
	conf.Passwd = dbConfig.DBPassword
	conf.DBName = dbConfig.DBName

	tlsConfig, err := dbConfig.TlsConfig(cnf.AppEnv)
	if err != nil {
		return nil, err
	}

	if tlsConfig != nil {
		if err = mysql.RegisterTLSConfig("custom", tlsConfig); err != nil {
			return nil, err
		}

		conf.TLSConfig = "custom"
	}

	return conf, nil
}

func ConnectDatabase(cnf *config.Config) (*sqlx.DB, error) {
	db, err := connect(cnf)

	if err != nil {
		return nil, err
	}
	return db, nil
}
