package httpserver

import (
	"fmt"
	"testLo/internal/domain"
	"time"
)

type taskReadDTO struct {
	ID        int64  `json:"id"`
	Title     string `json:"title"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type createTaskRequest struct {
	Title  string `json:"title"`
	Status string `json:"status"`
}

func toTaskReadDTO(t *domain.Task) taskReadDTO {
	statusString, err := getStringFromStatusForTask(t.Status)
	if err != nil {
		panicMsg := fmt.Sprintf("cannot convert task status %s to a string", t.Status)
		panic(panicMsg)
	}
	return taskReadDTO{
		ID:        t.ID,
		Title:     t.Title,
		Status:    statusString,
		CreatedAt: t.CreatedAt.Format(time.RFC3339),
		UpdatedAt: t.UpdatedAt.Format(time.RFC3339),
	}
}

func toRepositoryTask(req createTaskRequest) (domain.Task, error) {
	status, err := GetStatusFromStringForTask(req.Status)
	if err != nil {
		return domain.Task{}, fmt.Errorf("invalid task status: %s", req.Status)
	}
	return domain.Task{
		Title:     req.Title,
		Status:    status,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}
