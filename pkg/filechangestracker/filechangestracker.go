package filechangestracker

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"sync"
	"time"

	"github.com/danielboakye/filechangestracker/pkg/config"
	"github.com/osquery/osquery-go"
)

type FileChangesTracker struct {
	trackerLogger          *slog.Logger
	appLogger              *slog.Logger
	config                 *config.Config
	timerLastHeartbeat     time.Time
	mu                     sync.Mutex
	osqueryClient          *osquery.ExtensionManagerClient
	LogMutex               sync.Mutex
	lastProcessedTimestamp int64
}

func New(trackerLogger *slog.Logger, appLogger *slog.Logger, cfg *config.Config) *FileChangesTracker {
	return &FileChangesTracker{
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
	query := fmt.Sprintf("SELECT * FROM file_events WHERE target_path LIKE '%s%%'  AND time > %d;", f.config.Directory, f.lastProcessedTimestamp)
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

		changeTime, err := strconv.ParseInt(row["time"], 10, 64)
		if err != nil {
			continue
		}
		if changeTime > f.lastProcessedTimestamp {
			f.lastProcessedTimestamp = changeTime
		}
	}

	return nil
}

func (f *FileChangesTracker) IsTimerThreadAlive() bool {
	checkFrequency := time.Duration(f.config.CheckFrequency) * time.Second
	buffer := 30 * time.Second
	deadline := checkFrequency + buffer

	f.mu.Lock()
	defer f.mu.Unlock()

	return time.Since(f.timerLastHeartbeat) < deadline
}
