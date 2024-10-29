package bootstrap

import (
	"fmt"
	"open_url_service/pkg/config"
	"open_url_service/pkg/database/mysql"
	"open_url_service/pkg/logger"
)

func RegistryDatabase(cfg *config.Config) *mysql.DB {
	// Remove this code below if no need database
	db, err := mysql.ConnectDatabase(cfg)
	if err != nil {
		logger.Fatal(fmt.Sprintf("Failed to connect to database: %v", err))
	}

	return mysql.New(db, false, cfg.DatabaseConfig.DBName)
}
