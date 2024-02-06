package api

import (
	"net/http"
	"path/filepath"

	"github.com/joberly/demo-temporal/workflows"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.temporal.io/api/serviceerror"
	"go.temporal.io/sdk/client"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// ApiParams holds the dependencies for the API.
type ApiParams struct {
	fx.In
	Router *gin.Engine
	Logger *zap.Logger
	Config *Config
	Client client.Client
}

// Api is the API server.
type Api struct {
	router *gin.Engine
	logger *zap.Logger
	config *Config
	client client.Client
}

func New(params *ApiParams) (*Api, error) {
	return &Api{
		router: params.Router,
		logger: params.Logger,
		config: params.Config,
		client: params.Client,
	}, nil
}

func (a *Api) Run() {
	a.router.POST("/upload", a.uploadHandler)
	a.router.GET("/status/:id", a.statusHandler)
	a.router.GET("/health", a.healthHandler)
	a.router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	if err := a.router.Run(":8080"); err != nil {
		a.logger.Error("failed to start server", zap.Error(err))
	}
}

func (a *Api) uploadHandler(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		a.logger.Error("failed to parse form", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
		return
	}

	uuid := uuid.New().String()
	uploadFilePath := filepath.Join(a.config.UploadDir, uuid)

	if err := c.SaveUploadedFile(file, uploadFilePath); err != nil {
		a.logger.Error("failed to save file", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save file"})
		return
	}

	a.logger.Info("recieved file",
		zap.String("uuid", uuid),
		zap.String("path", uploadFilePath),
	)

	// start ImageProcessingWorkflow
	wfOpts := client.StartWorkflowOptions{
		ID:        uuid,
		TaskQueue: a.config.TaskQueue,
	}
	wfRun, err := a.client.ExecuteWorkflow(c.Request.Context(),
		wfOpts, workflows.ImageProcessingWorkflow, uuid)
	if err != nil {
		a.logger.Error("failed to start workflow", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to start process"})
		return
	}

	// workflow started successfully
	c.JSON(http.StatusAccepted,
		gin.H{
			"message":    "file uploaded",
			"imageId":    uuid,
			"workflowId": wfRun.GetID(),
			"runId":      wfRun.GetRunID(),
		},
	)
}

func (a *Api) statusHandler(c *gin.Context) {
	runID := c.Param("runId")
	workflowID := c.Param("workflowId")

	a.logger.Info("received image status request", zap.String("runId", runID))

	// query Temporal for the status of the workflow
	desc, err := a.client.DescribeWorkflowExecution(c.Request.Context(),
		workflowID, runID)
	if err != nil {
		switch err.(type) {
		case *serviceerror.InvalidArgument:
			a.logger.Info("workflow not found",
				zap.String("workflowId", workflowID),
				zap.String("runId", runID),
			)
			c.JSON(http.StatusNotFound, gin.H{
				"workflowId": workflowID,
				"runId":      runID,
				"error":      "not found",
			})
		default:
			a.logger.Error("failed to get status", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get status"})
		}
		return
	}

	c.JSON(http.StatusOK,
		gin.H{
			"workflowId": workflowID,
			"runId":      runID,
			"status":     desc.WorkflowExecutionInfo.GetStatus().String(),
		},
	)
}

func (a *Api) healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
