package statuses

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// Handle request for `GET /v1/statuses/{id}`
func (h *handler) FetchStatusByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// statusのidをURLパラメタから取得
	ids := chi.URLParam(r, "id")
	if ids == "" {
		http.Error(w, fmt.Sprintf("empty id is invalid to find a status"), http.StatusBadRequest)
		return
	}
	id, err := strconv.ParseInt(ids, 10 /* base */, 64 /* bitSize */)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// dbからidに対応したstatusを取得
	status, err := h.sr.FindStatusByID(ctx, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(status); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
