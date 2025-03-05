package tests

import (
	"net/http"
	"testing"

	config "github.com/DogFox/CutIt/configs"
)

func TestImageFromCache(t *testing.T) {
	config, err := config.NewConfig("../configs/config.yaml")
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	previewerURL := config.Tests.Url
	NginxURL := config.Tests.Nginx
	url := NginxURL + "test.jpg"

	t.Run("TestGetImageFromCache", func(t *testing.T) {
		resp, err := http.Get(url)
		if err != nil || resp.StatusCode != http.StatusOK {
			t.Fatalf("Expected 200, got %d, err: %v", resp.StatusCode, err)
		}
	})

	t.Run("TestRemoteServerNotExists", func(t *testing.T) {
		url := previewerURL + "http://nonexistent.domain/image.jpg"
		resp, err := http.Get(url)
		if err != nil || resp.StatusCode != http.StatusInternalServerError {
			t.Fatalf("Expected 500, got %d, err: %v", resp.StatusCode, err)
		}
	})

	t.Run("TestRemoteServerReturned404", func(t *testing.T) {
		url := previewerURL + "notfound.jpg"
		resp, err := http.Get(url)
		if err != nil || resp.StatusCode != http.StatusInternalServerError {
			t.Fatalf("Expected 500, got %d, err: %v", resp.StatusCode, err)
		}
	})

	t.Run("TestRemoteServerReturnedNonImage", func(t *testing.T) {
		url := previewerURL + "malware.exe"
		resp, err := http.Get(url)
		if err != nil || resp.StatusCode != http.StatusInternalServerError {
			t.Fatalf("Expected 500, got %d, err: %v", resp.StatusCode, err)
		}
	})

	t.Run("TestRemoteServerReturnedError", func(t *testing.T) {
		url := previewerURL + "error.jpg"
		resp, err := http.Get(url)
		if err != nil || resp.StatusCode != http.StatusInternalServerError {
			t.Fatalf("Expected 500, got %d, err: %v", resp.StatusCode, err)
		}
	})

	t.Run("TestRemoteServerReturnedImage", func(t *testing.T) {
		url := previewerURL + "test.jpg"
		resp, err := http.Get(url)
		if err != nil || resp.StatusCode != http.StatusInternalServerError {
			t.Fatalf("Expected 500, got %d, err: %v", resp.StatusCode, err)
		}
	})

	t.Run("TestImageSmallerThanRequiredSize", func(t *testing.T) {
		url := previewerURL + "small.jpg"
		resp, err := http.Get(url)
		if err != nil || resp.StatusCode != http.StatusInternalServerError {
			t.Fatalf("Expected 500, got %d, err: %v", resp.StatusCode, err)
		}
	})
}
