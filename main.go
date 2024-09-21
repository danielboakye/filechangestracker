package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/danielboakye/filechangestracker/internal/commandexecutor"
	"github.com/danielboakye/filechangestracker/pkg/config"
	"github.com/danielboakye/filechangestracker/pkg/filechangestracker"
	"github.com/danielboakye/filechangestracker/pkg/httpserver"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("error loading config: %w", err)
	}

	trackerLogFile, err := os.OpenFile(config.FileChangesLogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Error opening log file: %v", err)
	}
	defer trackerLogFile.Close()

	appLogger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	executor := commandexecutor.New(appLogger, cfg)
	if err := executor.Start(ctx); err != nil {
		log.Fatal("failed to start command executor: %w", err)
	}

	trackerLogger := slog.New(slog.NewJSONHandler(trackerLogFile, nil))
	tracker := filechangestracker.New(trackerLogger, appLogger, cfg)
	if err := tracker.Start(ctx); err != nil {
		log.Fatal("failed to start tracker: %w", err)
	}
	appLogger.Info("started-tracker-on-directory", slog.String("directory", cfg.Directory))

	addr := fmt.Sprintf(":%s", cfg.HTTPPort)
	apiServer := httpserver.NewServer(addr, appLogger, tracker, executor)
	if err := apiServer.Start(); err != nil {
		log.Fatal("failed to start http server on: ", addr)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-quit
	appLogger.Info("starting-shutdown")
	defer cancel()
	apiServer.Stop(ctx)
	executor.Stop(ctx)
	tracker.Stop(ctx)
	appLogger.Info("shutdown-complete!")

}
