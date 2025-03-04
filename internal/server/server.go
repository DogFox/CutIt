package internalhttp

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/DogFox/CutIt/internal/app"
	"github.com/DogFox/CutIt/internal/cache"
	"github.com/DogFox/CutIt/internal/logger"
)

type Server struct {
	Addr    string
	Handler http.Handler
	logg    *logger.Logger
	cache   *cache.Cache
	app     *app.App
}

type Logger interface {
	*logger.Logger
}

func NewServer(logger *logger.Logger, app *app.App, dsn string) *Server {
	mux := http.NewServeMux()
	server := &Server{
		logg: logger,
		app:  app,
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

	width := parts[2]
	height := parts[3]
	imgURL := strings.Join(parts[4:], "/")

	resizedPath, err := s.app.Resize(imgURL, width, height)
	if err != nil {
		http.Error(w, "Error processing image", http.StatusInternalServerError)
		return
	}

	http.ServeFile(w, r, resizedPath)
}
