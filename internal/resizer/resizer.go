package resizer

import (
	"log"
	"strconv"

	"github.com/disintegration/imaging"
)

func Resize(inputPath, outputPath, width, height string) error {
	img, err := imaging.Open(inputPath)
	if err != nil {
		return err
	}

	widthInt, err := strconv.Atoi(width)
	if err != nil {
		log.Println("Failed to cast str to int:", err)
	}
	heightInt, err := strconv.Atoi(height)
	if err != nil {
		log.Println("Failed to cast str to int:", err)
	}

	resized := imaging.Resize(img, widthInt, heightInt, imaging.Lanczos)

	err = imaging.Save(resized, outputPath)
	if err != nil {
		log.Println("Failed to save resized image:", err)
	}
	return err
}
