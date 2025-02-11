package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type AddRequest struct {
	UserId      int64  `json:"user_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	DueDate     string `json:"due_date"`
}

type AddResponse struct {
	TaskID int64  `json:"task_id"`
	Error  string `json:"error,omitempty"`
}

type Adder interface {
	AddTask(userId int64, taskName string, taskDescription string, dueDate time.Time) (int64, error)
}

func Add(adder Adder) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.task.add.New"

		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req AddRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, fmt.Sprintf("%s: invalid request body: %v", op, err), http.StatusBadRequest)
			return
		}

		dueDate, err := time.Parse(time.RFC3339, req.DueDate)
		if err != nil {
			http.Error(w, fmt.Sprintf("%s: invalid due_date format, use RFC3339", op), http.StatusBadRequest)
			return
		}

		taskID, err := adder.AddTask(req.UserId, req.Name, req.Description, dueDate)
		if err != nil {
			http.Error(w, fmt.Sprintf("%s: failed to add task: %v", op, err), http.StatusInternalServerError)
			return
		}

		resp := AddResponse{TaskID: taskID}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	}
}
