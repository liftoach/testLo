package httpserver

import (
	"fmt"
	"testLo/internal/domain"
)

var statusTaskMapping = map[string]domain.TaskStatus{
	"TaskDone":       domain.TaskStatusDone,
	"TaskInProgress": domain.TaskStatusInProgress,
	"TaskToDo":       domain.TaskStatusTodo,
}

func GetStatusFromStringForTask(status string) (domain.TaskStatus, error) {
	mappedStatus, ok := statusTaskMapping[status]
	if !ok {
		return 0, fmt.Errorf("invalid task status: %s", status)
	}
	return mappedStatus, nil
}

var taskStatusMapping = map[domain.TaskStatus]string{
	domain.TaskStatusDone:       "TaskDone",
	domain.TaskStatusInProgress: "TaskInProgress",
	domain.TaskStatusTodo:       "TaskToDo",
}

func getStringFromStatusForTask(status domain.TaskStatus) (string, error) {
	mappedStatus, ok := taskStatusMapping[status]
	if !ok {
		return "", fmt.Errorf("invalid task status: %w", status)
	}

	return mappedStatus, nil
}
