package app

import (
	"path/filepath"

	"github.com/DogFox/CutIt/internal/cache"
	"github.com/DogFox/CutIt/internal/downloader"
	"github.com/DogFox/CutIt/internal/resizer"
)

type App struct {
	logger Logger
	cache  *cache.Cache
}

type Logger interface{}

func New(logger Logger, cache *cache.Cache) *App {
	return &App{
		logger: logger,
		cache:  cache,
	}
}

func (a *App) Resize(imgURL, width, height string) (string, error) {
	cacheKey := filepath.Join("cache", width, height, imgURL)

	if path, found := a.cache.Get(cacheKey); found {
		return path, nil
	}

	originalPath := filepath.Join("cache", "original.jpg")
	resizedPath := filepath.Join("cache", "resized.jpg")

	if err := downloader.Download(imgURL, originalPath); err != nil {
		return "", err
	}

	if err := resizer.Resize(originalPath, resizedPath, width, height); err != nil {
		return "", err
	}

	a.cache.Put(cacheKey, resizedPath)
	return resizedPath, nil
}
