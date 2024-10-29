package config

import (
	"errors"
	"github.com/fsnotify/fsnotify"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"net/http"
	"open_url_service/pkg/logger"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	AppConfig      `mapstructure:",squash"`
	LoggerConfig   `mapstructure:",squash"`
	DatabaseConfig `mapstructure:",squash"`
	NoSleepConfig  `mapstructure:",squash"`
}

func LoadAllConfigs() (*Config, error) {
	var cnf Config
	err := loadConfig()
	if err != nil {
		return nil, err
	}

	err = viper.Unmarshal(&cnf)
	if err != nil {
		return nil, err
	}

	return &cnf, nil
}

// FiberConfig func for configuration Fiber app.
func (cnf *Config) FiberConfig() fiber.Config {
	// Return Fiber configuration.
	return fiber.Config{
		AppName:       cnf.AppName,
		StrictRouting: false,
		CaseSensitive: false,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			// Status code defaults to 500
			code := fiber.StatusInternalServerError

			var e *fiber.Error
			if errors.As(err, &e) {
				code = e.Code
			}

			switch code {
			// Not found handle
			case http.StatusNotFound:
				err := c.Status(http.StatusNotFound).JSON(map[string]interface{}{
					"message": "Sorry, the resource not found!",
				})
				if err != nil {
					return err
				}
			// Method not allowed handle
			case http.StatusMethodNotAllowed:
				err := c.Status(http.StatusMethodNotAllowed).JSON(map[string]interface{}{
					"message": "Method not allowed!",
				})
				if err != nil {
					return err
				}
			// Default internal server error handle
			case http.StatusInternalServerError:
				err := c.Status(http.StatusInternalServerError).JSON(map[string]interface{}{
					"message": "Something went wrong!",
				})
				if err != nil {
					return err
				}
			}
			return nil
		},
	}
}

func loadConfig() error {
	files, err := getAllConfigFiles("./config")
	if err != nil {
		logger.Warn(err)
	}

	viper.AddConfigPath("./config")
	for _, file := range files {
		viper.SetConfigType("json")
		viper.SetConfigFile(file)
		err = viper.MergeInConfig()
		if err != nil {
			return err
		}
	}

	viper.AutomaticEnv()

	viper.AddConfigPath(".")
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	_ = viper.MergeInConfig()

	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		logrus.Infof("Config %v was change", e.Name)
	})

	return nil
}

func getAllConfigFiles(folderPath string) ([]string, error) {
	var configFiles []string

	err := filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(info.Name(), ".json") {
			configFiles = append(configFiles, path)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return configFiles, nil
}
