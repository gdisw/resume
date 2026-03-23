package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Health struct{}

func (Health) MountOn(router chi.Router) {
	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
}
