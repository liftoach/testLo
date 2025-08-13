package httpserver

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"testLo/internal/domain"
	"testLo/internal/service"
)

type TaskHandler struct {
	service service.Service
}

func NewTaskHandler(s service.Service) *TaskHandler {
	return &TaskHandler{service: s}
}

// GET /tasks?status=todo
func (h *TaskHandler) GetTasks(w http.ResponseWriter, r *http.Request) {
	statusFilter := r.URL.Query().Get("status")
	var tasks []domain.Task
	var err error

	if statusFilter != "" {
		tasks, err = h.service.GetAllByStatus(r.Context(), statusFilter)
	} else {
		tasks, err = h.service.GetAll(r.Context())
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
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(toTaskReadDTO(task))
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

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(toTaskReadDTO(created))
}
