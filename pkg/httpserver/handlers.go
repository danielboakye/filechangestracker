package httpserver

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

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

	err = s.executor.AddCommands(req.Commands)
	if err != nil {
		response.InternalError(w)
	}

	response.JSON(w, http.StatusOK, map[string]string{
		"message": "commands added to queue",
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
	res, err := s.tracker.GetLogs()
	if err != nil {
		response.InternalError(w)
		return
	}

	response.JSON(w, http.StatusOK, res)
}

func (s *Server) NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	response.JSON(w, http.StatusNotFound, map[string]string{
		"message": fmt.Sprintf("resource: (%s) could not be found", r.URL.Path),
	})
}
