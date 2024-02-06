package activities

import (
	"context"
	"path/filepath"

	"go.uber.org/zap"
)

// CopyImageActivity is a Temporal activity that copies an image
// from the uploads location to the working location.
func (a *Activities) CopyImageActivity(ctx context.Context, imageID string) error {
	a.logger.Info("copying image", zap.String("imageID", imageID))

	// copy file from upload to working dir
	err := a.copyFile(filepath.Join(a.config.UploadDir, imageID),
		filepath.Join(a.config.WorkingDir, imageID))
	if err != nil {
		return err
	}

	a.logger.Info("image copied", zap.String("imageID", imageID))
	return nil
}
