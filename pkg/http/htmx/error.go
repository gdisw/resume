package htmx

import "errors"

var (
	ErrNoHtmxRequest = errors.New("htmx: not an HTMX request")
)
