package repository

import (
	"context"
	"fmt"
	"sync"
	"testLo/internal/domain"
	"testLo/internal/service"
	"time"
)

type MemoryTaskRepo struct {
	mu    sync.RWMutex
	tasks map[int64]domain.Task
	idSeq int64
	logCh chan string
}

func NewMemoryTaskRepo(logCh chan string) *MemoryTaskRepo {
	return &MemoryTaskRepo{
		tasks: make(map[int64]domain.Task),
		logCh: logCh,
	}
}

func (r *MemoryTaskRepo) log(msg string) {
	select {
	case r.logCh <- "[repository] " + msg:
	default:
	}
}

func (r *MemoryTaskRepo) GetAllTasks(ctx context.Context) ([]domain.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	select {
	case <-ctx.Done():
		r.log("GetAllTasks canceled by context")
		return nil, ctx.Err()
	default:
	}

	res := make([]domain.Task, 0, len(r.tasks))
	for _, t := range r.tasks {
		res = append(res, t)
	}
	return res, nil
}

func (r *MemoryTaskRepo) GetAllTasksByStatus(ctx context.Context, status string) ([]domain.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	select {
	case <-ctx.Done():
		r.log("GetAllTasksByStatus canceled by context")
		return nil, ctx.Err()
	default:
	}

	taskStatus, err := getStatusFromStringForTask(status)
	if err != nil {
		r.log(fmt.Sprintf("get status frin string for task: %v", err))
		return nil, nil
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

	select {
	case <-ctx.Done():
		r.log(fmt.Sprintf("GetTaskByID canceled by context (id=%d)", id))
		return nil, ctx.Err()
	default:
	}

	t, ok := r.tasks[id]
	if !ok {
		return nil, service.ErrNotFound
	}
	return &t, nil
}

func (r *MemoryTaskRepo) CreateTask(ctx context.Context, t domain.Task) (*domain.Task, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	select {
	case <-ctx.Done():
		r.log("CreateTask canceled by context")
		return nil, ctx.Err()
	default:
	}

	r.idSeq++
	t.ID = r.idSeq
	t.CreatedAt = time.Now()
	r.tasks[t.ID] = t

	return &t, nil
}
