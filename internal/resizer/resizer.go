package resizer

import (
	"strconv"

	"github.com/DogFox/CutIt/internal/logger"
	"github.com/disintegration/imaging"
)

type ImageCutter struct {
	logg *logger.Logger
}

func NewImageCutter(logg *logger.Logger) *ImageCutter {
	return &ImageCutter{
		logg: logg,
	}
}
func (i ImageCutter) Resize(inputPath, outputPath, width, height string) error {
	img, err := imaging.Open(inputPath)
	if err != nil {
		i.logg.Error("Unable to open input image")
		return err
	}

	widthInt, err := strconv.Atoi(width)
	if err != nil {
		i.logg.Error("Failed to cast str to int:")
		return err
	}
	heightInt, err := strconv.Atoi(height)
	if err != nil {
		i.logg.Error("Failed to cast str to int:")
		return err
	}

	resized := imaging.Resize(img, widthInt, heightInt, imaging.Lanczos)

	err = imaging.Save(resized, outputPath)
	if err != nil {
		i.logg.Error("Failed to save resized image")
		return err
	}
	return nil
}
