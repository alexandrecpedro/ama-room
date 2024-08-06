package api

import (
	"net/http"

	"github.com/alexandrecpedro/ama-room/backend/internal/store/pgstore"

	"github.com/go-chi/chi/v5"
)

type apiHandler struct {
	//? TODO: use an interface
	q *pgstore.Queries
	r *chi.Mux
}

func (h apiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.r.ServeHTTP(w, r)
}

func NewHandler(q *pgstore.Queries) http.Handler {
	a := apiHandler{
		q: q,
	}

	r := chi.NewRouter()

	a.r = r
	return a
}
