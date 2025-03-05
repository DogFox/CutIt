package downloader

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/DogFox/CutIt/internal/logger"
)

type Downloader struct {
	logg *logger.Logger
}

func NewDownloader(logg *logger.Logger) *Downloader {
	return &Downloader{
		logg: logg,
	}
}

func (t Downloader) Download(url, filename string, headers map[string]string) error {
	client := &http.Client{}

	req, err := http.NewRequest("GET", "http://"+url, nil)
	if err != nil {
		t.logg.Error("failed to fetch image: %w", err)
		return fmt.Errorf("failed to fetch image: %w", err)
	}

	// req.Header.Set("User-Agent", "ImagePreviewer/1.0")
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	resp, err := client.Do(req)
	if err != nil {
		t.logg.Error("failed to fetch image: %w", err)
		return fmt.Errorf("failed to fetch image: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.logg.Errorf("failed to fetch image, status: %d", resp.StatusCode)
		return fmt.Errorf("failed to fetch image, status: %d", resp.StatusCode)
	}

	out, err := os.Create(filename)
	if err != nil {
		t.logg.Errorf("failed to create file: %v", err)
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		t.logg.Errorf("failed to save image: %v", err)
		return fmt.Errorf("failed to save image: %w", err)
	}

	return nil
}
