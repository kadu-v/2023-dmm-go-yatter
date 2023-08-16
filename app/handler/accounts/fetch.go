package accounts

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type FetchRequest struct {
	ID int64
}

// Handle request for `GET /v1/accounts/{username}`
func (h *handler) Fetch(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// usernameをURLのパラメタから取得
	username := chi.URLParam(r, "username")
	if username == "" {
		http.Error(w, fmt.Sprintf("empty username is invalid to find a account"), http.StatusBadRequest)
		return
	}

	// Qusetion: ユーザーが見つからない状況はエラーなのか．それともnullを返してしまっていいのか？
	account, err := h.ar.FindByUsername(ctx, username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else if account == nil {
		http.Error(w, fmt.Sprintf("Not find a user \"%s\"", username), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(account); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
