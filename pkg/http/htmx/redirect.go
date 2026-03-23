package htmx

import (
	"net/http"
	"net/url"
)

func Redirect(w http.ResponseWriter, path string) {
	w.Header().Set("HX-Redirect", path)
}

func Refresh(w http.ResponseWriter, r *http.Request) {
	current := CurrentURL(r)
	Redirect(w, current.String())
}

func PushURL(w http.ResponseWriter, uri url.URL) {
	w.Header().Set("Hx-Push-Url", uri.String())
}
