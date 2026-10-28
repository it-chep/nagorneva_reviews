package middleware

import (
	"context"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"net/http"
	"time"
)

// Timeout cancels slow requests without writing a second status after the
// handler has already responded. Handlers must honor the request context.
func Timeout(duration time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, cancel := context.WithTimeout(r.Context(), duration)
			defer cancel()
			writer := chimiddleware.NewWrapResponseWriter(w, r.ProtoMajor)
			next.ServeHTTP(writer, r.WithContext(ctx))
			if ctx.Err() == context.DeadlineExceeded && writer.Status() == 0 {
				writer.WriteHeader(http.StatusGatewayTimeout)
			}
		})
	}
}
