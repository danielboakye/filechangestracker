package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/danielboakye/filechangestracker/internal/commandexecutor"
	"github.com/danielboakye/filechangestracker/internal/config"
	"github.com/danielboakye/filechangestracker/internal/filechangestracker"
	"github.com/danielboakye/filechangestracker/internal/httpserver"
	"github.com/danielboakye/filechangestracker/internal/mongolog"
	"github.com/danielboakye/filechangestracker/pkg/osquerymanager"
	"github.com/osquery/osquery-go"
)

func main() {
	cfg, err := config.LoadConfig(config.ConfigName, config.ConfigPath)
	if err != nil {
		log.Fatalf("error loading config: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logStore, err := mongolog.NewMongoLogStore(ctx, cfg.MongoURI, config.LogsDBName, config.LogsCollectionName)
	if err != nil {
		log.Fatalf("failed to start mongo: %v", err)
	}

	appLogger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	executor := commandexecutor.New(appLogger, cfg)
	if err := executor.Start(ctx); err != nil {
		log.Fatalf("failed to start command executor: %v", err)
	}

	osqueryClient, err := osquery.NewClient(cfg.SocketPath, 10*time.Second)
	if err != nil {
		log.Fatalf("Error creating osquery client: %v", err)
	}
	osqueryManager := osquerymanager.New(osqueryClient)

	tracker := filechangestracker.New(appLogger, cfg, osqueryManager, logStore)
	if err := tracker.Start(ctx); err != nil {
		log.Fatalf("failed to start tracker: %v", err)
	}
	appLogger.Info("started-tracker-on-directory", slog.String("directory", cfg.Directory))

	handler := httpserver.NewHandler(tracker, executor)
	router := handler.RegisterRoutes()

	addr := fmt.Sprintf(":%s", cfg.HTTPPort)
	apiServer := httpserver.NewServer(addr, appLogger, router)
	if err := apiServer.Start(); err != nil {
		log.Fatal("failed to start http server on: ", addr)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-quit

	appLogger.Info("starting-shutdown")
	apiServer.Stop(ctx)
	executor.Stop(ctx)
	tracker.Stop(ctx)
	logStore.Close(ctx)
	appLogger.Info("shutdown-complete!")
}
