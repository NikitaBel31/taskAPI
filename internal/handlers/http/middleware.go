package http

import (
	"context"
	"net/http"
	"taskapi/internal/logger"
	"time"
)

type ctxKey string

const requestIDKey ctxKey = "req_id"

func requestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := newReqID()
		ctx := context.WithValue(r.Context(), requestIDKey, reqID)
		w.Header().Set("X-Request-ID", reqID)
		ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func newReqID() string {
	return time.Now().UTC().Format("20060102T150405.000000000")
}

func requestIDFromCtx(ctx context.Context) string {
	v, _ := ctx.Value(requestIDKey).(string)
	return v
}

func (rt *Router) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		ww := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(ww, r)

		reqID := requestIDFromCtx(r.Context())
		rt.log.Log(logger.Entry{
			Time:      start.UTC(),
			Event:     "http_request",
			RequestID: reqID,
			Data: map[string]any{
				"method": r.Method,
				"path":   r.URL.Path,
				"status": ww.statusCode,
				"took":   time.Since(start).String(),
			},
		})
	})
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
