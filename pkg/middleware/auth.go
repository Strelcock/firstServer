package middleware

import (
	"context"
	"firstServer/configs"
	"firstServer/pkg/jwt"
	"net/http"
	"strings"
)

type key string

const (
	contextEmailKey key = "contextEmailKey"
)

func IsAuthed(next http.Handler, config *configs.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authedHeader := r.Header.Get("Authorization")
		token := strings.TrimPrefix(authedHeader, "Bearer ")
		_, data := jwt.NewJWT(config.Auth.Secret).Parse(token)
		r.Context()
		ctx := context.WithValue(r.Context(), contextEmailKey, data.Email)
		req := r.WithContext(ctx)
		next.ServeHTTP(w, req)
	})
}
