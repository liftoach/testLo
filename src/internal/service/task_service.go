package service

import (
	"context"
	"errors"
	"fmt"
	"testLo/internal/domain"
)

type TaskRepository interface {
	GetAllTasksByStatus(ctx context.Context, status string) ([]domain.Task, error)
	GetTaskByID(ctx context.Context, id int64) (*domain.Task, error)
	CreateTask(ctx context.Context, t domain.Task) (*domain.Task, error)
	GetAllTasks(ctx context.Context) ([]domain.Task, error)
}

func (s *Service) GetAllTasks(ctx context.Context) ([]domain.Task, error) {
	tasks, err := s.TaskRepository.GetAllTasks(ctx)
	if err != nil {
		s.log(fmt.Sprintf("task repository: get all tasks: %v", err))
	}
	return tasks, err
}

func (s *Service) GetAllTasksByStatus(ctx context.Context, status string) ([]domain.Task, error) {
	tasks, err := s.TaskRepository.GetAllTasksByStatus(ctx, status)
	if err != nil {
		s.log(fmt.Sprintf("task repository: get all tasks by status: %v", err))
	}
	return tasks, err
}

func (s *Service) GetTaskByID(ctx context.Context, id int64) (*domain.Task, error) {
	task, err := s.TaskRepository.GetTaskByID(ctx, id)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil, ErrNotFound
		}
		s.log(fmt.Sprintf("task repository: get all tasks by id: %v", err))
	}
	return task, nil
}

func (s *Service) CreateTask(ctx context.Context, t domain.Task) (*domain.Task, error) {
	task, err := s.TaskRepository.CreateTask(ctx, t)
	if err != nil {
		s.log(fmt.Sprintf("task repository: create task: %v", err))
	}
	return task, err
}
