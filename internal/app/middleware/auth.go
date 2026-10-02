package middleware

import (
	"net/http"
	"strings"
)

// TokenVerifier is the authentication dependency required by HTTP middleware.
type TokenVerifier interface {
	Verify(string) (int64, error)
}

// AdminAuth protects HTTP grpc-gateway endpoints that require an admin token.
// Authorization is deliberately bypassed when debug is enabled.
func AdminAuth(tokens TokenVerifier, debug bool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if debug {
				next.ServeHTTP(w, r)
				return
			}
			authorization := r.Header.Get("Authorization")
			if !strings.HasPrefix(authorization, "Bearer ") {
				http.Error(w, `{"error":"authorization required"}`, http.StatusUnauthorized)
				return
			}
			if _, err := tokens.Verify(strings.TrimPrefix(authorization, "Bearer ")); err != nil {
				http.Error(w, `{"error":"invalid token"}`, http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
