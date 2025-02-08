package handlers

import "to-do-list/internal/task"

type ShowHandler interface {
	ShowTask(taskName string) (task.Task, error)
}
