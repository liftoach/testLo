package service

import (
	"context"
	"testLo/internal/domain"
)

type TaskRepository interface {
	GetAllTasksByStatus(ctx context.Context, status string) ([]domain.Task, error)
	GetTaskByID(ctx context.Context, id int64) (*domain.Task, error)
	CreateTask(ctx context.Context, t domain.Task) (*domain.Task, error)
	GetAllTasks(ctx context.Context) ([]domain.Task, error)
}

func (s *Service) GetAll(ctx context.Context) ([]domain.Task, error) {
	return s.TaskRepository.GetAllTasks(ctx)
}

func (s *Service) GetAllByStatus(ctx context.Context, status string) ([]domain.Task, error) {
	return s.TaskRepository.GetAllTasksByStatus(ctx, status)
}

func (s *Service) GetTaskByID(ctx context.Context, id int64) (*domain.Task, error) {
	task, err := s.TaskRepository.GetTaskByID(ctx, id)
	if err == nil {
		return nil, err
	}
	return task, err
}

func (s *Service) CreateTask(ctx context.Context, t domain.Task) (*domain.Task, error) {
	task, err := s.TaskRepository.CreateTask(ctx, t)
	if err == nil {
		return nil, err
	}
	return task, err
}
