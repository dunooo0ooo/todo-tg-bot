package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type ChangeRequest struct {
	TaskId      int64  `json:"task_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	DueDate     string `json:"due_date"`
}

type ChangeResponse struct {
	Error string `json:"error,omitempty"`
}

type Changer interface {
	ChangeTask(taskId int64, newName string, newDescription string, newDueDate time.Time) error
}

func ChangeHandler(ch Changer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.Change.task"

		if r.Method != http.MethodPut {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req ChangeRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, fmt.Sprintf("%s: invalid request body: %v", op, err), http.StatusBadRequest)
			return
		}

		newDueDate, err := time.Parse("2006-01-02", req.DueDate)
		if err != nil {
			http.Error(w, fmt.Sprintf("%s: invalid due_date format (expected YYYY-MM-DD): %v", op, err), http.StatusBadRequest)
			return
		}

		err = ch.ChangeTask(req.TaskId, req.Name, req.Description, newDueDate)
		if err != nil {
			http.Error(w, fmt.Sprintf("%s: failed to change task: %v", op, err), http.StatusInternalServerError)
			return
		}

		resp := ChangeResponse{}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	}
}
