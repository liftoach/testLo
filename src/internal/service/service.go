package service

import (
	"errors"
)

var (
	ErrNotFound = errors.New("not found")
)

type Service struct {
	TaskRepository TaskRepository
	logCh          chan string
}

func NewService(taskRepo TaskRepository, logCh chan string) *Service {
	return &Service{
		TaskRepository: taskRepo,
		logCh:          logCh,
	}
}

func (s *Service) log(msg string) {
	select {
	case s.logCh <- "[service] " + msg:
	default:
	}
}
