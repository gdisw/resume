package htmx

import "net/http"

const (
	SwapInnerHTML   = SwapType("innerHTML")
	SwapOuterHTML   = SwapType("outerHTML")
	SwapBeforeBegin = SwapType("beforebegin")
	SwapAfterBegin  = SwapType("afterbegin")
	SwapBeforeEnd   = SwapType("beforeend")
	SwapAfterEnd    = SwapType("afterend")
	SwapDelete      = SwapType("delete")
	SwapNone        = SwapType("none")
)

type SwapType string

func Reswap(w http.ResponseWriter, st SwapType) {
	w.Header().Set("Hx-Reswap", string(st))
}

func Retarget(w http.ResponseWriter, target string) {
	w.Header().Set("Hx-Retarget", target)
}
