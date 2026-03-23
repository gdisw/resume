package csrf

import (
	"net/http"
	"regexp"

	"github.com/justinas/nosurf"
)

const (
	exemptTypeFunction = exemptType(iota)
	exemptTypeGlob
	exemptTypePath
	exemptTypeRegex
)

type exemptType uint8
type exemptCfg struct {
	exType exemptType
	s      string
	fn     func(*http.Request) bool
	r      *regexp.Regexp
}

func ExemptFunc(fn func(*http.Request) bool) exemptCfg {
	return exemptCfg{
		exType: exemptTypeFunction,
		fn:     fn,
	}
}

func ExemptGlob(glob string) exemptCfg {
	return exemptCfg{
		exType: exemptTypeGlob,
		s:      glob,
	}
}

func ExemptPath(path string) exemptCfg {
	return exemptCfg{
		exType: exemptTypePath,
		s:      path,
	}
}

func ExemptRegexp(r *regexp.Regexp) exemptCfg {
	return exemptCfg{
		exType: exemptTypePath,
		r:      r,
	}
}

func WithExempts(cfgs ...exemptCfg) handlerOption {
	return func(h *nosurf.CSRFHandler) {
		var t exemptType

		for _, cfg := range cfgs {
			switch t {
			case exemptTypeFunction:
				h.ExemptFunc(cfg.fn)
			case exemptTypeGlob:
				h.ExemptGlob(cfg.s)
			case exemptTypePath:
				h.ExemptPath(cfg.s)
			case exemptTypeRegex:
				h.ExemptRegexp(cfg.r)
			}
		}
	}
}
