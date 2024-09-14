package bootstrap

import (
	"open_url_service/pkg/config"
	"open_url_service/pkg/logger"
	"open_url_service/pkg/util"
)

func RegistryLogger(cfg *config.Config) {
	loggerConfig := logger.Config{
		Environment: util.EnvironmentTransform(cfg.AppEnv),
		Debug:       cfg.AppDebug,
		Level:       cfg.LogLevel,
		ServiceName: cfg.AppName,
	}

	logger.Setup(loggerConfig)

	switch cfg.LogDriver {
	case logger.LogDriverLoki:
		logger.AddHook(cfg.LoggerConfig.WithLokiHook(cfg))
		break
	default:
		break
	}
}
