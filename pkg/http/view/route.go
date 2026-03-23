package view

import (
	"context"
	"net/http"
	"path"
	"strings"
)

type RouteConfig struct {
	PathPrefix string
	Exclude    []string
}

func (parent RouteConfig) push(child RouteConfig) RouteConfig {
	parent.PathPrefix = path.Join(parent.PathPrefix, child.PathPrefix)
	return parent
}

func route(viewData ViewData, p string) string {
	raw, ok := viewData["RouteConfig"]
	if !ok {
		return p
	}

	cfg, ok := raw.(RouteConfig)
	if !ok {
		return p
	}

	for _, pattern := range cfg.Exclude {
		if strings.HasPrefix(p, pattern) {
			return p
		}
	}

	return path.Join(cfg.PathPrefix, p)
}

func PutRouteConfig(cfg RouteConfig) func(http.Handler) http.Handler {
	if !strings.HasPrefix(cfg.PathPrefix, "/") {
		panic("a route config path prefix must start with '/'")
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			root, ok := r.Context().Value(contextKeyRouteConfig).(RouteConfig)
			if ok {
				root = root.push(cfg)
			} else {
				root = cfg
			}

			next.ServeHTTP(w, r.WithContext(
				ContextWithViewData(
					context.WithValue(r.Context(), contextKeyRouteConfig, root),
					ViewData{"RouteConfig": root},
				),
			))
		})
	}
}
