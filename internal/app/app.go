package app

import (
	"fmt"
	"net/url"
	"os"
	"path"
	"path/filepath"

	"github.com/DogFox/CutIt/internal/cache"
	"github.com/DogFox/CutIt/internal/downloader"
	"github.com/DogFox/CutIt/internal/logger"
	"github.com/DogFox/CutIt/internal/resizer"
)

type App struct {
	logger     *logger.Logger
	cache      *cache.Cache
	downloader *downloader.Downloader
	cutter     *resizer.ImageCutter
}

type Logger interface{}

func New(logger *logger.Logger, cache *cache.Cache, downloader *downloader.Downloader, imageCutter *resizer.ImageCutter) *App {
	return &App{
		logger:     logger,
		cache:      cache,
		downloader: downloader,
		cutter:     imageCutter,
	}

}

func (a *App) Resize(imgURL, width, height string, headers map[string]string) (string, error) {
	parsedURL, err := url.Parse(imgURL)
	if err != nil {
		return "", fmt.Errorf("invalid URL: %s", imgURL)
	}
	fileName := path.Base(parsedURL.Path)
	cacheKey := filepath.Join("cache", width, height, fileName)

	if path, found := a.cache.Get(cacheKey); found {
		return path, nil
	}

	originalPath := filepath.Join("cache", fileName)
	resizedPath := filepath.Join("cache", width+"_"+height+"_"+fileName)

	if err := os.MkdirAll(filepath.Dir(originalPath), os.ModePerm); err != nil {
		return "", fmt.Errorf("unable to make dir")
	}

	if err := a.downloader.Download(imgURL, originalPath, headers); err != nil {
		return "", fmt.Errorf("unable to download image")
	}

	if err := a.cutter.Resize(originalPath, resizedPath, width, height); err != nil {
		return "", fmt.Errorf("unable to resize image")
	}

	a.cache.Put(cacheKey, resizedPath)
	return resizedPath, nil
}
