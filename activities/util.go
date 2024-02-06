package activities

import (
	"io"
	"os"

	"go.uber.org/zap"
)

func (a *Activities) copyFile(src, dst string) error {
	// open source file in uploads dir
	srcFile, err := os.Open(src)
	if err != nil {
		a.logger.Error("failed to open source file", zap.Error(err))
		return err
	}
	defer srcFile.Close()

	// create destination file in working dir
	dstFile, err := os.Create(dst)
	if err != nil {
		a.logger.Error("failed to create destination file", zap.Error(err))
		return err
	}
	defer dstFile.Close()

	// copy file contents
	if _, err := io.Copy(dstFile, srcFile); err != nil {
		a.logger.Error("failed to copy file contents", zap.Error(err))
		return err
	}

	return nil
}
