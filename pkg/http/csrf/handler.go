package csrf

import (
	"net/http"
	"regexp"

	"github.com/justinas/nosurf"
)

type handlerOption func(*nosurf.CSRFHandler)

func Handler(next http.Handler) http.Handler {
	handler := nosurf.New(next)
	handler.ExemptRegexp(regexp.MustCompile("^/track/.*"))

	// nosurf generates a cookie on each route with the path
	// set to the request path. We want the CSRF tokens to be
	// valid on all paths and therefore set the base cookie
	// path to be the root (i.e. /).
	handler.SetBaseCookie(http.Cookie{Path: "/"})

	return handler
}
