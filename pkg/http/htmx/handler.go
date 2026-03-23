package htmx

import (
	"net/http"
	"net/url"
)

func Check(r *http.Request) bool {
	return r.Header.Get("HX-Request") == "true"
}

func CurrentURL(r *http.Request) url.URL {
	uri, err := url.Parse(r.Header.Get("Hx-Current-Url"))
	if err == nil {
		return *uri
	}
	return url.URL{}
}

func Target(r *http.Request) string {
	return r.Header.Get("Hx-Target")
}

func TriggerId(r *http.Request) string {
	return r.Header.Get("Hx-Trigger")
}
