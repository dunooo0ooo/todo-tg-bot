package handlers

type DeleteHandler interface {
	DeleteTask(taskName string) error
}
