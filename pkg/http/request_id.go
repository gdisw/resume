package http

import (
	"context"
	"net/http"

	"github.com/gdisw/resume/pkg/identifier"
	"github.com/go-chi/chi/v5/middleware"
)

const RequestPrefix = "req"

func RequestId(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rid := identifier.Mint(RequestPrefix)
		ctx := context.WithValue(r.Context(), middleware.RequestIDKey, rid)

		w.Header().Set(middleware.RequestIDHeader, rid)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
