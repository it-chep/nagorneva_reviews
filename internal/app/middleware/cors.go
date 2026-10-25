package middleware

import "net/http"

// CORSMiddleware permits browser requests only from explicitly configured
// origins. Requests without Origin are server-to-server and remain allowed.
func CORSMiddleware(corsConfig OriginConfig) func(next http.Handler) http.Handler {
	allowedOrigins := allowedOriginSet(corsConfig.GetAllowedOrigins())

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := normalizeOrigin(r.Header.Get("Origin"))
			_, isAllowedOrigin := allowedOrigins[origin]
			if origin != "" && isAllowedOrigin {
				w.Header().Set("Vary", "Origin")
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
				w.Header().Set("Access-Control-Allow-Credentials", "true")
			}

			if r.Method == http.MethodOptions {
				if origin == "" || !isAllowedOrigin {
					w.WriteHeader(http.StatusForbidden)
					return
				}
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
