package api

import (
	"encoding/json"

	"github.com/spf13/viper"
	"go.temporal.io/sdk/client"
	"go.uber.org/zap"
)

// Config holds the configuration for the API.
type Config struct {
	UploadDir    string
	ProcessedDir string
	TemporalHost string
	TemporalPort string
	TaskQueue    string
}

func NewConfig(logger *zap.Logger) (*Config, error) {
	viper.AutomaticEnv()
	viper.SetEnvPrefix("demo")
	viper.SetDefault("UPLOAD_DIR", "/tmp/uploads")
	viper.SetDefault("PROCESSED_DIR", "/tmp/processed")
	viper.SetDefault("TEMPORAL_HOST", "localhost")
	viper.SetDefault("TEMPORAL_PORT", "7233")
	viper.SetDefault("TASK_QUEUE", "image-processing")

	config := &Config{
		UploadDir:    viper.GetString("UPLOAD_DIR"),
		ProcessedDir: viper.GetString("PROCESSED_DIR"),
		TemporalHost: viper.GetString("TEMPORAL_HOST"),
		TemporalPort: viper.GetString("TEMPORAL_PORT"),
		TaskQueue:    viper.GetString("TASK_QUEUE"),
	}

	configJson, err := json.Marshal(config)
	if err != nil {
		logger.Error("failed to marshal configuration", zap.Error(err))
		return nil, err
	}

	logger.Info("loaded configuration", zap.String("config", string(configJson)))
	return config, nil
}

func NewTemporalClient(config *Config, logger *zap.Logger) (client.Client, error) {
	return client.Dial(client.Options{
		HostPort: config.TemporalHost + ":" + config.TemporalPort,
	})
}
