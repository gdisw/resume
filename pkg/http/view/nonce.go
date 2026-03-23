package view

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"net/http"
)

func ScriptNonceFromContext(ctx context.Context) string {
	nonce, _ := ctx.Value(contextKeyScriptNonce).(string)
	return nonce
}

func PutScriptNonce(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nonce := genScriptNonce()
		next.ServeHTTP(w, r.WithContext(
			ContextWithViewData(
				context.WithValue(r.Context(), contextKeyScriptNonce, nonce),
				ViewData{"ScriptNonce": nonce},
			),
		))
	})
}

func genScriptNonce() string {
	buf := make([]byte, 16)
	_, _ = rand.Read(buf)
	return base64.
		StdEncoding.
		WithPadding(base64.NoPadding).
		EncodeToString(buf)
}
