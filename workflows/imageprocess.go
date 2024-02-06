package workflows

import (
	"github.com/joberly/demo-temporal/activities"

	"go.temporal.io/sdk/workflow"
)

// ImageProcessingWorkflow is a Temporal workflow that processes an image.
func ImageProcessingWorkflow(ctx workflow.Context, imageID string) error {
	// This is where you would put your image processing logic.
	// For now, we'll just log the image ID.
	workflow.GetLogger(ctx).Info("Processing image", "imageID", imageID)

	// Run an activity to check the image size.
	// This is a blocking call, so the workflow will wait for the result.
	// The result will be stored in the `size` variable.
	var size int
	err := workflow.ExecuteActivity(ctx, activities.CheckImageSizeActivity, imageID).Get(ctx, &size)
	if err != nil {
		// Return the error to indicate that the workflow failed.
		return err
	}

	// Return nil to indicate that the workflow completed successfully.
	return nil
}
