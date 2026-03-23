package csrf

import (
	"net/http"

	"github.com/justinas/nosurf"
)

func CookiePathClearer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if path != "" && path != "/" {
			http.SetCookie(w, &http.Cookie{
				Name:   nosurf.CookieName,
				Value:  "",
				Path:   path,
				MaxAge: -1,
			})
		}

		next.ServeHTTP(w, r)
	})
}
