package handlers

import "time"

type ChangeHandler interface {
	ChangeTask(oldName string, newName string, newDescription string, newDueDate time.Time) error
}
