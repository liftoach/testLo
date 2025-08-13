package service

import (
	"errors"
)

var (
	ErrNotFound = errors.New("not found")
)

type Service struct {
	TaskRepository TaskRepository
}

func NewService(taskRepo TaskRepository) *Service {
	return &Service{
		TaskRepository: taskRepo,
	}
}
