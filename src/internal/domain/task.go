package domain

import (
	"time"
)

type TaskStatus uint8

const (
	TaskStatusDone TaskStatus = 1 + iota
	TaskStatusInProgress
	TaskStatusTodo
)

type Task struct {
	ID        int64
	Title     string
	Status    TaskStatus
	CreatedAt time.Time
	UpdatedAt time.Time
}
