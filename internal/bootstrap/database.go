package bootstrap

import (
	"fmt"
	"open_url_service/pkg/config"
	"open_url_service/pkg/database/postgres"
	"open_url_service/pkg/logger"
)

func RegistryDatabase(cfg *config.Config) *postgres.DB {
	// Remove this code below if no need database
	db, err := postgres.ConnectDatabase(cfg)
	if err != nil {
		logger.Fatal(fmt.Sprintf("Failed to connect to database: %v", err))
	}

	return postgres.New(db, false, cfg.DatabaseConfig.DBName)
}
