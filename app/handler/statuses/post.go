package statuses

import (
	"encoding/json"
	"net/http"
	"yatter-backend-go/app/domain/object"
	"yatter-backend-go/app/handler/auth"
)

// Request body for `POST /v1/statuses`
type PostRequest struct {
	Status string `json:"status"`
}

// Handle request for `POST /v1/statuses`
func (h *handler) PostStatus(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req PostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Status objectの各フィールドの初期化
	status := new(object.Status)
	account := auth.AccountOf(r)
	status.Account = account
	status.Content = req.Status

	// db へstatusをpost
	if ID, err := h.sr.AddStatus(ctx, account, status); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else {
		status.ID = ID
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(status); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
