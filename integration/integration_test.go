package tests

import (
	"net/http"
	"testing"
	"time"
)

const previewerURL = "http://localhost:8050/fill/300/200/"
const nginxURL = "localhost:8081/images/"

func waitForService(url string) bool {
	timeout := time.After(10 * time.Second)
	tick := time.Tick(500 * time.Millisecond)

	for {
		select {
		case <-timeout:
			return false
		case <-tick:
			resp, err := http.Get(url)
			if err == nil && resp.StatusCode < 500 {
				return true
			}
		}
	}
}

func TestImageFromCache(t *testing.T) {
	url := previewerURL + nginxURL + "test.jpg"
	// if !waitForService(url) {
	// 	t.Fatal("Service not available")
	// }

	resp, err := http.Get(url)
	if err != nil || resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected 200, got %d, err: %v", resp.StatusCode, err)
	}
}

func TestRemoteServerNotExists(t *testing.T) {
	url := previewerURL + "http://nonexistent.domain/image.jpg"
	resp, err := http.Get(url)
	if err != nil || resp.StatusCode != http.StatusBadGateway {
		t.Fatalf("Expected 502, got %d, err: %v", resp.StatusCode, err)
	}
}

func TestRemoteServerReturned404(t *testing.T) {
	url := previewerURL + "notfound.jpg"
	resp, err := http.Get(url)
	if err != nil || resp.StatusCode != http.StatusNotFound {
		t.Fatalf("Expected 404, got %d, err: %v", resp.StatusCode, err)
	}
}

func TestRemoteServerReturnedNonImage(t *testing.T) {
	url := previewerURL + "malware.exe"
	resp, err := http.Get(url)
	if err != nil || resp.StatusCode != http.StatusUnsupportedMediaType {
		t.Fatalf("Expected 415, got %d, err: %v", resp.StatusCode, err)
	}
}

func TestRemoteServerReturnedImage(t *testing.T) {
	url := previewerURL + "test.jpg"
	resp, err := http.Get(url)
	if err != nil || resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected 200, got %d, err: %v", resp.StatusCode, err)
	}
}

func TestImageSmallerThanRequiredSize(t *testing.T) {
	url := previewerURL + "small.jpg"
	resp, err := http.Get(url)
	if err != nil || resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected 200, got %d, err: %v", resp.StatusCode, err)
	}
}
