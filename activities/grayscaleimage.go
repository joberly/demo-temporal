package activities

import (
	"context"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"

	"golang.org/x/image/webp"

	"go.uber.org/zap"
)

// GrayscaleImageActivity is a Temporal activity that converts a working image
// to black and white.
func (a *Activities) GrayscaleImageActivity(ctx context.Context, imageID string) error {
	a.logger.Info("converting image to grayscale", zap.String("imageID", imageID))

	// open image file
	a.logger.Info("opening working image", zap.String("imageID", imageID))
	file, err := os.Open(filepath.Join(a.config.WorkingDir, imageID))
	if err != nil {
		return err
	}

	// decode for image format
	a.logger.Info("decoding working image type", zap.String("imageID", imageID))
	_, format, err := image.DecodeConfig(file)
	if err != nil {
		file.Close()
		return err
	}

	// rewind file to beginning
	a.logger.Info("rewinding working image", zap.String("imageID", imageID))
	_, err = file.Seek(0, 0)
	if err != nil {
		file.Close()
		return err
	}

	// decode image
	a.logger.Info("decoding working image data", zap.String("imageID", imageID))
	var img image.Image
	switch format {
	case "jpeg":
		img, err = jpeg.Decode(file)
	case "png":
		img, err = png.Decode(file)
	case "webp":
		img, err = webp.Decode(file)
	default:
		err = fmt.Errorf("unsupported image format: %s", format)
	}
	// done with the source image file at this point
	file.Close()
	if err != nil {
		return err
	}

	// convert the image to grayscale
	a.logger.Info("converting image to grayscale", zap.String("imageID", imageID))
	gray := a.convertToGrayscale(ctx, img)

	// save the grayscale image as a jpeg in the processed dir
	a.logger.Info("creating processed image file", zap.String("imageID", imageID))
	processedFile, err := os.Create(filepath.Join(a.config.ProcessedDir, imageID))
	if err != nil {
		return err
	}
	defer processedFile.Close()

	// encode the grayscale image as jpeg
	a.logger.Info("writing grayscale image file", zap.String("imageID", imageID))
	err = jpeg.Encode(processedFile, gray, &jpeg.Options{Quality: 90})
	if err != nil {
		return err
	}

	a.logger.Info("conversion to grayscale complete", zap.String("imageID", imageID))
	return nil
}

func (a *Activities) convertToGrayscale(ctx context.Context, img image.Image) *image.Gray {
	bounds := img.Bounds()
	grayImg := image.NewGray(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			originalColor := img.At(x, y)
			grayColor := color.GrayModel.Convert(originalColor)
			grayImg.Set(x, y, grayColor)
		}
	}
	return grayImg
}
