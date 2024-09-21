package httpserver

import (
	"bufio"
	"encoding/json"
	"io"
	"net/http"
	"os"

	"github.com/danielboakye/filechangestracker/pkg/config"
	"github.com/danielboakye/filechangestracker/pkg/response"
)

// CommandRequest represents the structure of a command request
type CommandRequest struct {
	Commands []string `json:"commands"`
}

// HealthCheckResponse represents the structure of the health check response
type HealthCheckResponse struct {
	WorkerThread bool `json:"worker_thread_alive"`
	TimerThread  bool `json:"timer_thread_alive"`
}

// LogsResponse represents the structure of logs response
type LogsResponse struct {
	Logs []string `json:"logs"`
}

func (s *Server) HandleSubmitCommands(w http.ResponseWriter, r *http.Request) {
	var req CommandRequest
	body, err := io.ReadAll(r.Body)
	if err != nil {
		response.InvalidRequest(w, err.Error())
		return
	}

	err = json.Unmarshal(body, &req)
	if err != nil {
		response.InvalidRequest(w, err.Error())
		return
	}
	if len(req.Commands) == 0 {
		response.InvalidRequest(w, "no commands submitted")
		return
	}

	s.executor.AddCommands(req.Commands)

	response.JSON(w, http.StatusOK, map[string]string{
		"message": "Commands added to queue",
	})
}

// handleHealthCheck returns the health status of the worker and timer threads
func (s *Server) HandleHealthCheck(w http.ResponseWriter, r *http.Request) {
	res := HealthCheckResponse{
		WorkerThread: s.executor.IsWorkerThreadAlive(),
		TimerThread:  s.tracker.IsTimerThreadAlive(),
	}

	response.JSON(w, http.StatusOK, res)
}

func (s *Server) HandleGetLogs(w http.ResponseWriter, r *http.Request) {
	s.tracker.LogMutex.Lock()
	defer s.tracker.LogMutex.Unlock()

	file, err := os.Open(config.FileChangesLogFile)
	if err != nil {
		response.InternalError(w)
		return
	}
	defer file.Close()

	var res []map[string]interface{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var jsonObject map[string]interface{}
		line := scanner.Text()

		err := json.Unmarshal([]byte(line), &jsonObject)
		if err != nil {
			response.InternalError(w)
			return
		}

		res = append(res, jsonObject)
	}

	if err := scanner.Err(); err != nil {
		response.InternalError(w)
		return
	}

	response.JSON(w, http.StatusOK, res)
}

func (s *Server) NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	response.JSON(w, http.StatusNotFound, map[string]string{
		"message": "The requested resource could not be found",
	})
}
