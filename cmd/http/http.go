package http

import (
	"fmt"
	"open_url_service/pkg/app"
	"open_url_service/pkg/config"
	"open_url_service/pkg/logger"
)

func Start() {
	logger.SetJSONFormatter()
	cnf, err := config.LoadAllConfigs()
	if err != nil {
		logger.Fatal(fmt.Sprintf("Failed to load configuration file: %v", err))
	}

	app.InitializeApp(cnf)
	application := app.GetServer()

	//go func() {
	//	StartServiceListener(cnf)
	//}()

	if err := application.StartServer(); err != nil {
		logger.Fatal(fmt.Sprintf("Failed to start server: %v", err))
	}

}
