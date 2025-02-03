package middleware

import (
	"context"
	"github.com/google/uuid"
	"net/http"
)

func ContextRequestMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := r.Header.Get("X-Request-ID")
		if reqID == "" {
			reqID = uuid.New().String()
		}

		ctx := context.WithValue(r.Context(), "requestID", reqID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
