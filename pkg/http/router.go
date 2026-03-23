package http

import (
	"net/http"
	"time"

	"github.com/gdisw/resume/pkg/http/session"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Router struct {
	mux      *chi.Mux
	sessions session.Store
	mws      []func(http.Handler) http.Handler
}

type Routes interface {
	MountOn(chi.Router)
}

func NewRouter(sessions session.Store) *Router {
	mux := chi.NewRouter()
	mux.Use(RequestId)
	mux.Use(middleware.Logger)
	mux.Use(DefaultContentType("text/html"))
	mux.Use(middleware.Compress(5, "text/html"))
	mux.Use(middleware.Recoverer)
	mux.Use(middleware.Timeout(45 * time.Second))

	mux.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	return &Router{mux: mux, sessions: sessions}
}

func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	router.mux.ServeHTTP(w, r)
}

func (router *Router) Attach(r Routes, mws ...func(http.Handler) http.Handler) {
	r.MountOn(router.mux.With(mws...))
}

func (router *Router) AttachProtected(r Routes) {
	r.MountOn(router.mux.With(router.mws...))
}

func (router *Router) Plug(mw func(http.Handler) http.Handler) {
	router.mws = append(router.mws, mw)
}
