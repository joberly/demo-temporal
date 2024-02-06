package api

import (
	"encoding/json"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type Config struct {
	ImageDir string
}

func NewConfig(logger *zap.Logger) (*Config, error) {
	viper.AutomaticEnv()
	viper.SetEnvPrefix("demo")
	viper.SetDefault("IMAGE_DIR", "/tmp")

	config := &Config{
		ImageDir: viper.GetString("IMAGE_DIR"),
	}

	configJson, err := json.Marshal(config)
	if err != nil {
		logger.Error("Failed to marshal configuration", zap.Error(err))
		return nil, err
	}

	logger.Info("Loaded configuration", zap.String("config", string(configJson)))
	return config, nil
}
