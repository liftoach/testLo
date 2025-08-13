package repository

import (
	"context"
	"sync"
	"testLo/internal/domain"
	"testLo/internal/service"
	"time"
)

type MemoryTaskRepo struct {
	mu    sync.RWMutex
	tasks map[int64]domain.Task
	idSeq int64
}

func NewMemoryTaskRepo() *MemoryTaskRepo {
	return &MemoryTaskRepo{
		tasks: make(map[int64]domain.Task),
	}
}

func (r *MemoryTaskRepo) GetAllTasks(ctx context.Context) ([]domain.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	res := make([]domain.Task, 0, len(r.tasks))
	for _, t := range r.tasks {
		res = append(res, t)
	}
	return res, nil
}

func (r *MemoryTaskRepo) GetAllTasksByStatus(ctx context.Context, status string) ([]domain.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	taskStatus, err := getStatusFromStringForTask(status)
	if err != nil {
		return nil, err
	}

	var res []domain.Task
	for _, t := range r.tasks {
		if t.Status == taskStatus {
			res = append(res, t)
		}
	}
	return res, nil
}

func (r *MemoryTaskRepo) GetTaskByID(ctx context.Context, id int64) (*domain.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	t, ok := r.tasks[id]
	if !ok {
		return nil, service.ErrNotFound
	}
	return &t, nil
}

func (r *MemoryTaskRepo) CreateTask(ctx context.Context, t domain.Task) (*domain.Task, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.idSeq++
	t.ID = r.idSeq
	t.CreatedAt = time.Now()
	r.tasks[t.ID] = t
	return &t, nil
}
