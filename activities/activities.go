package activities

import (
	"go.uber.org/zap"
)

type Config struct {
	UploadDir    string
	WorkingDir   string
	ProcessedDir string
}

type ActivitiesParams struct {
	Logger *zap.Logger
	Config *Config
}

type Activities struct {
	logger *zap.Logger
	config *Config
}

func New(p *ActivitiesParams) *Activities {
	return &Activities{
		logger: p.Logger,
		config: p.Config,
	}
}
