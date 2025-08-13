package httpserver

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"testLo/internal/domain"
	"testLo/internal/service"
)

type TaskHandler struct {
	service service.Service
	logCh   chan string
}

func NewTaskHandler(s service.Service, logCh chan string) *TaskHandler {
	return &TaskHandler{service: s, logCh: logCh}
}

func (h *TaskHandler) log(msg string) {
	select {
	case h.logCh <- "[handler] " + msg:
	default:
	}
}

func (h *TaskHandler) withRecovery(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				h.log(fmt.Sprintf("panic recovered: %v", rec))
				http.Error(w, "internal server error", http.StatusInternalServerError)
			}
		}()
		fn(w, r)
	}
}

func (h *TaskHandler) HandleRoutes(mux *http.ServeMux, logCh chan<- string) {
	mux.HandleFunc("/tasks", h.withRecovery(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			h.GetTasks(w, r)
			h.log("GET /tasks called")
		case http.MethodPost:
			h.CreateTask(w, r)
			h.log("POST /tasks called")
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	mux.HandleFunc("/tasks/", h.withRecovery(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			h.GetTaskByID(w, r)
			logCh <- "GET /tasks/{id} called"
		} else {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	}))
}

func (h *TaskHandler) GetTasks(w http.ResponseWriter, r *http.Request) {
	statusFilter := r.URL.Query().Get("status")
	var tasks []domain.Task
	var err error

	if statusFilter != "" {
		tasks, err = h.service.GetAllTasksByStatus(r.Context(), statusFilter)
		h.log("GET /tasks with status=" + statusFilter)
	} else {
		tasks, err = h.service.GetAllTasks(r.Context())
		h.log("GET /tasks")
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var dtos []taskReadDTO
	for _, t := range tasks {
		dtos = append(dtos, toTaskReadDTO(&t))
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(dtos)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// GET /tasks/{id}
func (h *TaskHandler) GetTaskByID(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/tasks/")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	task, err := h.service.GetTaskByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
		}
		h.log(fmt.Sprintf("service: get task by id: %v", err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(toTaskReadDTO(task)); err != nil {
		h.log(fmt.Sprintf("failed to encode task to JSON: %v", err))
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
}

// POST /tasks
func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var req createTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	task, err := toRepositoryTask(req)
	if err != nil {
		http.Error(w, "invalid status", http.StatusBadRequest)
		return
	}

	created, err := h.service.CreateTask(r.Context(), task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err = json.NewEncoder(w).Encode(toTaskReadDTO(created)); err != nil {
		h.log(fmt.Sprintf("failed to encode task to JSON: %v", err))
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
}
