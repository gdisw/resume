package htmx

import (
	"net/http"
)

func Dispatch(isHtmx, otherwise http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if Check(r) {
			isHtmx(w, r)
			return
		}

		otherwise(w, r)
	}
}
