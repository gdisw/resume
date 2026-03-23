package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Home struct{}

func (Home) MountOn(router chi.Router) {
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello, World!"))
	})
}
