package worker

import (
	"github.com/joberly/demo-temporal/activities"
	"github.com/joberly/demo-temporal/workflows"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type WorkerParams struct {
	fx.In
	Logger *zap.Logger
	Config *Config
	Client client.Client
}

type Worker struct {
	logger *zap.Logger
	config *Config
	client client.Client
	worker worker.Worker
}

func New(params WorkerParams) (*Worker, error) {
	return &Worker{
		logger: params.Logger,
		config: params.Config,
		client: params.Client,
	}, nil
}

func (w *Worker) Run() {
	// create a Temporal worker
	w.worker = worker.New(w.client, w.config.TaskQueue, worker.Options{})

	// register workflows
	w.worker.RegisterWorkflow(workflows.ImageProcessingWorkflow)

	// create activities and register them
	acts := activities.New(&activities.ActivitiesParams{
		Logger: w.logger,
		Config: &activities.Config{
			UploadDir:    w.config.UploadDir,
			WorkingDir:   w.config.WorkingDir,
			ProcessedDir: w.config.ProcessedDir,
		},
	})

	// register activities
	w.worker.RegisterActivity(acts.CopyImageActivity)
	w.worker.RegisterActivity(acts.GrayscaleImageActivity)

	// start the worker
	w.worker.Run(worker.InterruptCh())
}
