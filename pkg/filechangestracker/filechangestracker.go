package filechangestracker

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/danielboakye/filechangestracker/pkg/config"
	"github.com/danielboakye/filechangestracker/pkg/osquerymanager"
)

//go:generate mockgen -destination=../../mocks/filechangestracker/mock_filechangestracker.go -package=filechangestrackermock -source=filechangestracker.go
type FileChangesTracker interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error

	IsTimerThreadAlive() bool
	GetLogs() ([]map[string]interface{}, error)
}

type fileChangesTracker struct {
	trackerLogger          *slog.Logger
	appLogger              *slog.Logger
	config                 *config.Config
	timerLastHeartbeat     time.Time
	mu                     sync.Mutex
	osqueryManager         osquerymanager.OSQueryManager
	logMutex               sync.Mutex
	lastProcessedTimestamp int64
}

func New(trackerLogger *slog.Logger, appLogger *slog.Logger, cfg *config.Config, osqueryManager osquerymanager.OSQueryManager) FileChangesTracker {
	return &fileChangesTracker{
		trackerLogger:  trackerLogger,
		appLogger:      appLogger,
		config:         cfg,
		osqueryManager: osqueryManager,
	}
}

func (f *fileChangesTracker) Start(ctx context.Context) error {
	go f.timerThread(ctx)

	return nil
}

func (f *fileChangesTracker) Stop(ctx context.Context) error {
	f.osqueryManager.Close()
	return nil
}

func (f *fileChangesTracker) timerThread(ctx context.Context) {
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

func (f *fileChangesTracker) checkFileChanges() error {
	query := fmt.Sprintf("SELECT * FROM file_events WHERE target_path LIKE '%s%%'  AND time > %d;", f.config.Directory, f.lastProcessedTimestamp)
	res, err := f.osqueryManager.Query(query)
	if err != nil {
		if errors.Is(err, osquerymanager.ErrNoChangesFound) {
			return nil
		}
		return fmt.Errorf("error querying file changes: %w", err)
	}

	f.logMutex.Lock()
	defer f.logMutex.Unlock()
	for _, row := range res {
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

func (f *fileChangesTracker) IsTimerThreadAlive() bool {
	checkFrequency := time.Duration(f.config.CheckFrequency) * time.Second
	buffer := 30 * time.Second
	deadline := checkFrequency + buffer

	f.mu.Lock()
	defer f.mu.Unlock()

	return time.Since(f.timerLastHeartbeat) < deadline
}

func (f *fileChangesTracker) GetLogs() ([]map[string]interface{}, error) {
	f.logMutex.Lock()
	defer f.logMutex.Unlock()

	file, err := os.Open(f.config.LogFile)
	if err != nil {
		f.appLogger.Error("error-getting-logs", slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}
	defer file.Close()

	var res []map[string]interface{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var jsonObject map[string]interface{}
		line := scanner.Text()

		err := json.Unmarshal([]byte(line), &jsonObject)
		if err != nil {
			f.appLogger.Error("error-getting-logs", slog.String("error", err.Error()))
			return nil, fmt.Errorf("failed to Unmarshal logs: %w", err)
		}

		res = append(res, jsonObject)
	}

	if err := scanner.Err(); err != nil {
		f.appLogger.Error("error-getting-logs", slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to scan log file: %w", err)
	}

	return res, nil
}
