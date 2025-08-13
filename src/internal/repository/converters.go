package repository

import (
	"fmt"
	"testLo/internal/domain"
)

var statusTaskMapping = map[string]domain.TaskStatus{
	"TaskDone":       domain.TaskStatusDone,
	"TaskInProgress": domain.TaskStatusInProgress,
	"TaskToDo":       domain.TaskStatusTodo,
}

func getStatusFromStringForTask(status string) (domain.TaskStatus, error) {
	mappedStatus, ok := statusTaskMapping[status]
	if !ok {
		return 0, fmt.Errorf("invalid task status: %s", status)
	}
	return mappedStatus, nil
}
