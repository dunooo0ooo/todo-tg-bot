package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"to-do-list/internal/task"
)

type ShowRequest struct {
	UserId int64 `json:"user_id"`
}

type ShowResponse struct {
	Tasks []task.Task `json:"tasks"`
}

type Shower interface {
	GetTasks(userId int64) ([]task.Task, error)
}

func ShowTasks(shower Shower) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.show.task"

		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req ShowRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, fmt.Sprintf("%s: invalid request body: %v", op, err), http.StatusBadRequest)
			return
		}

		tasks, err := shower.GetTasks(req.UserId)
		if err != nil {
			http.Error(w, fmt.Sprintf("%s: error getting tasks: %v", op, err), http.StatusInternalServerError)
			return
		}

		resp := ShowResponse{Tasks: tasks}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, fmt.Sprintf("%s: error encoding response: %v", op, err), http.StatusInternalServerError)
			return
		}
	}
}
