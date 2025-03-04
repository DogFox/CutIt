package main

import (
	"context"
	"flag"
	"fmt"
	"os/signal"
	"sync"
	"syscall"

	config "github.com/DogFox/CutIt/configs"
	app "github.com/DogFox/CutIt/internal/app"
	logger "github.com/DogFox/CutIt/internal/logger"
	http "github.com/DogFox/CutIt/internal/server"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "../../configs/config.yaml", "Path to configuration file")
}

// func logEvery10Seconds() {
// 	for {
// 		fmt.Println("Logging every 10 seconds...")
// 		time.Sleep(10 * time.Second)
// 	}
// }

// go run ./cmd/ -config ./configs/config.yaml.
func main() {
	flag.Parse()

	config, err := config.NewConfig(configFile)
	if err != nil {
		fmt.Println(configFile)
		fmt.Println(err)
		return
	}
	logg := logger.New(config.Logger.Level)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	imageCutter := app.New(logg)

	httpServer := http.NewServer(logg, imageCutter, config.Server.DSN())

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		startHTTPServer(ctx, httpServer, logg)
	}()

	wg.Wait()
	logg.Info("calendar is running...")
}

func startHTTPServer(ctx context.Context, server *http.Server, logger *logger.Logger) error {
	go func() {
		<-ctx.Done()
		if err := server.Stop(ctx); err != nil {
			logger.Error("failed to stop http server: " + err.Error())
		}
		logger.Error("server stopped")
	}()

	logger.Println("http server started: ", server.Addr)
	return server.Start(ctx)
}
