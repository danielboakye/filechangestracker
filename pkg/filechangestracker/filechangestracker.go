package filechangestracker

import (
	"context"
	"fmt"
	"log/slog"
	"os/exec"
	"sync"
	"time"

	"github.com/danielboakye/filechangestracker/pkg/config"
	"github.com/osquery/osquery-go"
)

type FileChangesTracker struct {
	commandQueue        chan string
	trackerLogger       *slog.Logger
	appLogger           *slog.Logger
	config              *config.Config
	workerLastHeartbeat time.Time
	timerLastHeartbeat  time.Time
	mu                  sync.Mutex
	osqueryClient       *osquery.ExtensionManagerClient
	LogMutex            sync.Mutex
}

func New(trackerLogger *slog.Logger, appLogger *slog.Logger, cfg *config.Config) *FileChangesTracker {
	return &FileChangesTracker{
		commandQueue:  make(chan string, 100),
		trackerLogger: trackerLogger,
		appLogger:     appLogger,
		config:        cfg,
	}
}

func (f *FileChangesTracker) Start(ctx context.Context) error {
	client, err := osquery.NewClient(f.config.SocketPath, 10*time.Second)
	if err != nil {
		return fmt.Errorf("Error creating osquery client: %w", err)
	}
	f.osqueryClient = client

	go f.workerThread(ctx)
	go f.timerThread(ctx)

	return nil
}

func (f *FileChangesTracker) Stop(ctx context.Context) error {
	if f.osqueryClient != nil {
		f.osqueryClient.Close()
		f.osqueryClient = nil
	}

	return nil
}

func (f *FileChangesTracker) workerThread(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Minute) // Heartbeat every minute
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			f.mu.Lock()
			f.workerLastHeartbeat = time.Now()
			f.mu.Unlock()
		case newCmd := <-f.commandQueue:
			err := f.executeCommand(newCmd)
			if err != nil {
				f.appLogger.Error("error-executing-command", slog.String("error", err.Error()))
			}
		}
	}
}

func (f *FileChangesTracker) executeCommand(command string) error {
	// cmd = exec.Command("cmd", "/C", command) // windows
	cmd := exec.Command("/bin/sh", "-c", command)
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("error executing command: %w", err)
	}

	return nil
}

func (f *FileChangesTracker) timerThread(ctx context.Context) {
	checkFrequency := time.Duration(f.config.CheckFrequency) * time.Second
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(checkFrequency):
			f.mu.Lock()
			f.timerLastHeartbeat = time.Now()
			f.mu.Unlock()

			err := f.checkFileChanges()
			if err != nil {
				f.appLogger.Error("error-checking-file-changes", slog.String("error", err.Error()))
			}
		}
	}
}

func (f *FileChangesTracker) checkFileChanges() error {
	query := "SELECT * FROM file_events WHERE target_path LIKE '" + f.config.Directory + "%';"
	res, err := f.osqueryClient.Query(query)
	if err != nil {
		return fmt.Errorf("error running osquery: %w", err)
	}
	if res.Status.Code != 0 {
		return fmt.Errorf("error running osquery: %s", res.Status.Message)
	}

	f.LogMutex.Lock()
	defer f.LogMutex.Unlock()
	for _, row := range res.Response {
		f.appLogger.Debug("new change detected", slog.String("target_path", row["target_path"]))
		f.trackerLogger.Info(
			"change detected",
			slog.Any("details", row),
		)
	}

	return nil
}

func (f *FileChangesTracker) IsWorkerThreadAlive() bool {
	f.mu.Lock()
	defer f.mu.Unlock()

	return time.Since(f.workerLastHeartbeat) < 2*time.Minute
}

func (f *FileChangesTracker) IsTimerThreadAlive() bool {
	checkFrequency := time.Duration(f.config.CheckFrequency) * time.Second
	buffer := 30 * time.Second
	deadline := checkFrequency + buffer

	f.mu.Lock()
	defer f.mu.Unlock()

	return time.Since(f.timerLastHeartbeat) < deadline
}

func (f *FileChangesTracker) AddCommands(commands []string) error {
	for _, cmd := range commands {
		f.commandQueue <- cmd
	}

	return nil
}
