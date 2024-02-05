package api

import (
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

type Api struct {
	Router *gin.Engine
	Logger *zap.Logger
	Config *Config
}

func New(router *gin.Engine, logger *zap.Logger, config *Config) (*Api, error) {
	return &Api{
		Router: router,
		Logger: logger,
		Config: config,
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
	file, err := c.FormFile("file")
	if err != nil {
		a.Logger.Error("Failed to parse form", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
		return
	}

	uuid := uuid.New().String()
	ingestFilePath := filepath.Join(a.Config.ImageDir, uuid+filepath.Ext(file.Filename))

	if err := c.SaveUploadedFile(file, ingestFilePath); err != nil {
		a.Logger.Error("Failed to save file", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save file"})
		return
	}

	a.Logger.Info("recieved file",
		zap.String("uuid", uuid),
		zap.String("path", ingestFilePath),
	)
	c.JSON(http.StatusAccepted, gin.H{"message": "file uploaded", "id": uuid})
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
