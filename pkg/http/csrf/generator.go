package csrf

import (
	"net/http"

	"github.com/gdisw/resume/pkg/http/htmx"
	"github.com/gdisw/resume/pkg/http/view"
	"github.com/justinas/nosurf"
)

func Generator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		if r.Method == http.MethodGet &&
			!htmx.Check(r) {
			ctx = view.ContextWithViewData(ctx, view.ViewData{
				"CSRF": nosurf.Token(r),
			})
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
