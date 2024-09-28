package filechangestracker

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"testing"
	"time"

	osquerymanagermock "github.com/danielboakye/filechangestracker/mocks/osquerymanager"
	"github.com/danielboakye/filechangestracker/pkg/config"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// go test -v -cover ./pkg/filechangestracker/...

// go test -v -cover -run TestCheckFileChanges ./pkg/filechangestracker
func TestCheckFileChanges(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	mockCtrl := gomock.NewController(t)
	mockOSQueryManager := osquerymanagermock.NewMockOSQueryManager(mockCtrl)

	cfg := &config.Config{
		LogFile: "test.log",
	}

	trackerLogFile, err := os.OpenFile(cfg.LogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	require.NoError(err)
	defer trackerLogFile.Close()

	appLogger := slog.Default()
	trackerLogger := slog.New(slog.NewJSONHandler(trackerLogFile, nil))

	tracker := New(trackerLogger, appLogger, cfg, mockOSQueryManager)
	it := tracker.(*fileChangesTracker)

	mockOSQueryManager.EXPECT().Query(gomock.Any()).Return([]map[string]string{
		{
			"target_path": "test/test.txt",
			"time":        strconv.FormatInt(time.Now().Unix(), 10),
		},
	}, nil).AnyTimes()

	err = it.checkFileChanges()
	assert.Nil(err)

	res, err := tracker.GetLogs()
	assert.Nil(err)
	assert.NotNil(res)
	assert.Len(res, 1)

	err = os.Remove(cfg.LogFile)
	require.NoError(err)
}

// go test -v -cover -run TestHealthCheck ./pkg/filechangestracker
func TestHealthCheck(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	mockCtrl := gomock.NewController(t)
	mockOSQueryManager := osquerymanagermock.NewMockOSQueryManager(mockCtrl)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg := &config.Config{
		LogFile:        "test.log",
		CheckFrequency: 1,
	}

	trackerLogFile, err := os.OpenFile(cfg.LogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	require.NoError(err)
	defer trackerLogFile.Close()

	appLogger := slog.Default()
	trackerLogger := slog.New(slog.NewJSONHandler(trackerLogFile, nil))

	tracker := New(trackerLogger, appLogger, cfg, mockOSQueryManager)

	mockOSQueryManager.EXPECT().Query(gomock.Any()).Return(nil, fmt.Errorf("no matches found")).AnyTimes()

	err = tracker.Start(ctx)
	require.NoError(err)

	time.Sleep(2 * time.Second)

	isAlive := tracker.IsTimerThreadAlive()
	assert.True(isAlive)

	res, err := tracker.GetLogs()
	assert.Nil(err)
	assert.Len(res, 0)

	err = os.Remove(cfg.LogFile)
	require.NoError(err)
}
