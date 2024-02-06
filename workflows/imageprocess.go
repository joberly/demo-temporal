package workflows

import (
	"time"

	"go.temporal.io/sdk/workflow"
)

// ImageProcessingWorkflowStatus is the status of an image processing workflow.
type ImageProcessingWorkflowStatus struct {
	ImageID string
	Status  string
	Error   string
}

// ImageProcessingWorkflow is a Temporal workflow that processes an image.
func ImageProcessingWorkflow(ctx workflow.Context, imageID string) error {
	workflow.GetLogger(ctx).Info("starting ImageProcessingWorkflow", "imageID", imageID)

	// setup a timeout for all activities
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute * 5,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	// status of the workflow reported back via query
	status := ImageProcessingWorkflowStatus{
		ImageID: imageID,
		Status:  "starting",
	}

	// setup a query handler to report the status of the workflow via the api
	err := workflow.SetQueryHandler(ctx, "status",
		func() (ImageProcessingWorkflowStatus, error) {
			return status, nil
		},
	)
	if err != nil {
		status.Status = "error"
		status.Error = err.Error()
		return err
	}

	workflow.GetLogger(ctx).Info("processing image", "imageID", imageID)

	// copy image to working directory
	status.Status = "copying image"
	err = workflow.ExecuteActivity(ctx, "CopyImageActivity", imageID).Get(ctx, nil)
	if err != nil {
		status.Status = "error copying image"
		status.Error = err.Error()
		return err
	}

	// convert image to grayscale
	status.Status = "converting image to grayscale"
	err = workflow.ExecuteActivity(ctx, "GrayscaleImageActivity", imageID).Get(ctx, nil)
	if err != nil {
		status.Status = "error converting image to grayscale"
		status.Error = err.Error()
		return err
	}

	// workflow successfully completed
	workflow.GetLogger(ctx).Info("image processing complete", "imageID", imageID)
	return nil
}
