package accounts

import (
	"encoding/json"
	"fmt"
	"net/http"

	"yatter-backend-go/app/domain/object"
)

// Request body for `POST /v1/accounts`
type AddRequest struct {
	Username string
	Password string
}

// Handle request for `POST /v1/accounts`
func (h *handler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req AddRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	account := new(object.Account)
	account.Username = req.Username
	if err := account.SetPassword(req.Password); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// アカウント名が1以上255以下かの確認
	usernameLength := len(account.Username)
	if !(1 <= usernameLength && usernameLength <= 255) {
		http.Error(w, fmt.Sprintf("invalid the lenght of username: %d", usernameLength), http.StatusBadRequest)
		return
	}

	// すでに同じアカウント名のユーザーが存在するか確認
	entity, err := h.ar.FindByUsername(ctx, account.Username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else if entity != nil {
		http.Error(w, fmt.Sprintf("user name \"%s\" is arleady existed", entity.Username), http.StatusBadRequest)
		return
	}

	// TODO: avator, header, note, created_atを埋めるべき
	// dbへ新規アカウントを登録
	if err := h.ar.CreateUser(ctx, account); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// dbからアカウントの作成日時を取得する．
	if entity, err := h.ar.FindByUsername(ctx, account.Username); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else if entity == nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else {
		account.CreateAt = entity.CreateAt
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(account); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
