package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

type Api struct {
	Router *gin.Engine
	Logger *zap.Logger
}

func New(router *gin.Engine, logger *zap.Logger) (*Api, error) {
	return &Api{
		Router: router,
		Logger: logger,
	}, nil
}

func (a *Api) Run() {
	a.Router.POST("/upload", a.uploadHandler)
	a.Router.GET("/status/:id", a.statusHandler)
	a.Router.GET("/health", a.healthHandler)
	a.Router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	if err := a.Router.Run(":8080"); err != nil {
		a.Logger.Error("Failed to start server", zap.Error(err))
	}
}

func (a *Api) uploadHandler(c *gin.Context) {
	// This is where you would put your image upload logic.
	// For now, we'll just log the image ID.
	a.Logger.Info("Received image upload request", zap.String("imageID", "123"))
	c.JSON(http.StatusAccepted, gin.H{"id": "123"})
}

func (a *Api) statusHandler(c *gin.Context) {
	imageID := c.Param("id")
	// This is where you would put your image status logic.
	// For now, we'll just log the image ID.
	a.Logger.Info("Received image status request", zap.String("imageID", imageID))
	c.JSON(http.StatusOK, gin.H{"status": "processing"})
}

func (a *Api) healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
