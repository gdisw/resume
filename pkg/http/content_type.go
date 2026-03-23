package http

import "net/http"

func DefaultContentType(ctype string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)

			header, name := w.Header(), "Content-Type"
			if header.Get(name) == "" {
				header.Set(name, ctype)
			}
		})
	}
}
