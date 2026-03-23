package view

import "context"

const (
	contextKeyViewData = contextKey(iota)
	contextKeyRouteConfig
	contextKeyScriptNonce
)

type contextKey uint8

func ContextWithViewData(ctx context.Context, viewData ViewData) context.Context {
	current, _ := ctx.Value(contextKeyViewData).([]ViewData)
	return context.WithValue(ctx, contextKeyViewData, append(current, viewData))
}

func viewDataFromContext(ctx context.Context) ViewData {
	data, _ := ctx.Value(contextKeyViewData).([]ViewData)
	merge := make(ViewData)
	for _, viewData := range data {
		for key, value := range viewData {
			merge[key] = value
		}
	}

	return merge
}
