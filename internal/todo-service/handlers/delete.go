package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type DeleteRequest struct {
	UserId int64  `json:"user_id"`
	Name   string `json:"name"`
}

type DeleteResponse struct {
	Error string `json:"error,omitempty"`
}

type Deleter interface {
	DeleteTask(userId int64, taskName string) error
}

func DeleteHandler(d Deleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.Delete.task"

		if r.Method != http.MethodDelete {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req DeleteRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, fmt.Sprintf("%s: invalid request body: %v", op, err), http.StatusBadRequest)
			return
		}

		var resp DeleteResponse
		err := d.DeleteTask(req.UserId, req.Name)
		if err != nil {
			resp.Error = err.Error()
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(resp)
			return
		}

		resp = DeleteResponse{Error: err.Error()}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	}
}
