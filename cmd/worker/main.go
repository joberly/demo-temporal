package main

import (
	"github.com/joberly/demo-temporal/internal/worker"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func main() {
	fx.New(
		fx.Provide(
			NewLogger,
			worker.NewConfig,
			worker.NewTemporalClient,
			worker.New,
		),
		fx.Invoke(func(w *worker.Worker) {
			w.Run()
		}),
	).Run()
}

func NewLogger() (*zap.Logger, error) {
	logger, err := zap.NewDevelopment()
	if err != nil {
		return nil, err
	}
	return logger, nil
}
