package main

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"net/http"
)

func TraceMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		traceID := uuid.New().String()
		ctx := context.WithValue(r.Context(), traceIDKey, traceID)
		logger := slog.Default().With("traceID", traceID)
		logger.InfoContext(ctx, "Request received", "path", r.URL.Path)	
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}