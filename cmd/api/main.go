package main

import (
	"github.com/joberly/demo-temporal/internal/api"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func main() {
	fx.New(
		fx.Provide(
			NewLogger,
			NewRouter,
			api.New,
		),
		fx.Invoke(func(a *api.Api) {
			a.Run()
		}),
	).Run()
}

func NewRouter() *gin.Engine {
	router := gin.Default()
	return router
}

func NewLogger() (*zap.Logger, error) {
	logger, err := zap.NewDevelopment()
	if err != nil {
		return nil, err
	}
	return logger, nil
}
