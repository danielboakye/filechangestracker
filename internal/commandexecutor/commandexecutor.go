package commandexecutor

import (
	"context"
	"fmt"
	"log/slog"
	"os/exec"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/danielboakye/filechangestracker/pkg/config"
)

type CommandExecutor interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error

	IsWorkerThreadAlive() bool
	AddCommands(commands []string) error
}

type commandExecutor struct {
	commandQueue        chan string
	appLogger           *slog.Logger
	config              *config.Config
	mu                  sync.Mutex
	workerLastHeartbeat time.Time
}

// List of commands to blacklist
var commandBlacklist = []string{
	"rm",       // Remove files or directories
	"shutdown", // Shutdown system
	"reboot",   // Reboot system
	"dd",       // Disk operations

	"mkfs",       // Format filesystem
	"del",        // Delete files
	"format",     // Format disk
	"powershell", // PowerShell execution
	"rd",         // Remove directory
}

func New(appLogger *slog.Logger, cfg *config.Config) CommandExecutor {
	return &commandExecutor{
		commandQueue: make(chan string, 100),
		appLogger:    appLogger,
		config:       cfg,
	}
}

func (f *commandExecutor) workerThread(ctx context.Context) {
	ticker := time.NewTicker(10 * time.Second) // Heartbeat every 10 seconds
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

// Check if a command is blacklisted
func isBlacklistedCommand(command string) bool {
	for _, cmdSubstr := range commandBlacklist {
		if strings.Contains(command, cmdSubstr) {
			return true
		}
	}
	return false
}

func (f *commandExecutor) executeCommand(command string) error {
	if isBlacklistedCommand(command) {
		return fmt.Errorf("execution blocked: blacklisted command detected: %s", command)
	}

	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/C", command) // Windows command execution
	} else {
		cmd = exec.Command("/bin/sh", "-c", command)
	}

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("error executing command: %w", err)
	}

	return nil
}

func (f *commandExecutor) Start(ctx context.Context) error {
	go f.workerThread(ctx)

	return nil
}

func (f *commandExecutor) Stop(ctx context.Context) error {
	return nil
}

func (f *commandExecutor) IsWorkerThreadAlive() bool {
	f.mu.Lock()
	defer f.mu.Unlock()

	return time.Since(f.workerLastHeartbeat) < 2*time.Minute
}

func (f *commandExecutor) AddCommands(commands []string) error {
	for _, cmd := range commands {
		f.commandQueue <- cmd
	}

	return nil
}
