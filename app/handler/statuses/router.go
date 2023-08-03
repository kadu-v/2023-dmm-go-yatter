package statuses

import (
	"net/http"
	"yatter-backend-go/app/domain/repository"
	"yatter-backend-go/app/handler/auth"

	"github.com/go-chi/chi/v5"
)

type handler struct {
	sr repository.Status
}

// Create Handler for `/v1/statuses/`
func NewRouter(ar repository.Account, sr repository.Status) http.Handler {
	r := chi.NewRouter()
	h := &handler{sr}

	// Statusのポスト
	r.Route("/", func(r chi.Router) {
		r.Use(auth.Middleware(ar))
		r.Post("/", h.PostStatus)
	})

	// 対応するidのstatusの取得
	r.Get("/{id}", h.FetchStatusByID)
	return r
}
