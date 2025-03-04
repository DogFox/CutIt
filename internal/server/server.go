package internalhttp

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/DogFox/CutIt/internal/cache"
	"github.com/DogFox/CutIt/internal/downloader"
	"github.com/DogFox/CutIt/internal/logger"
	"github.com/DogFox/CutIt/internal/resizer"
)

type Server struct {
	Addr    string
	Handler http.Handler
	logg    *logger.Logger
	cache   *cache.Cache
}

type Logger interface {
	*logger.Logger
}

type Application interface{}

func NewServer(logger *logger.Logger, app Application, dsn string) *Server {
	mux := http.NewServeMux()
	server := &Server{
		logg: logger,
	}
	mux.HandleFunc("/fill", server.ResizeImage)

	return server
}

func (s *Server) Start(ctx context.Context) error {
	server := &http.Server{
		Addr:         s.Addr,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      s.Handler,
	}

	s.logg.Infoln("start server on ", s.Addr)
	err := server.ListenAndServe()
	if err != nil {
		fmt.Println(err)
		s.logg.Errorln(err)
	}

	<-ctx.Done()
	return nil
}

func (s *Server) Stop(_ context.Context) error {
	os.Exit(1)
	return nil
}

func (s *Server) ResizeImage(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		http.Error(w, "Invalid URL format", http.StatusBadRequest)
		return
	}

	width, height := parts[2], parts[3]
	imgURL := strings.Join(parts[4:], "/")

	cacheKey := width + "x" + height + "_" + imgURL
	if path, found := s.cache.Get(cacheKey); found {
		http.ServeFile(w, r, path)
		return
	}

	originalPath := filepath.Join("cache", "original.jpg")
	resizedPath := filepath.Join("cache", "resized.jpg")

	if err := downloader.Download(imgURL, originalPath); err != nil {
		http.Error(w, "Failed to download image", http.StatusBadGateway)
		return
	}

	if err := resizer.Resize(originalPath, resizedPath, width, height); err != nil {
		http.Error(w, "Failed to resize image", http.StatusInternalServerError)
		return
	}

	s.cache.Put(cacheKey, resizedPath)
	http.ServeFile(w, r, resizedPath)
}
