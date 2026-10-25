package utils

import (
	"net/http"
	"os"
	"strings"
)

// CORSMiddleware adds CORS headers to every response.
// Allowed origins are read from the ALLOWED_ORIGINS env var (comma-separated).
// Falls back to localhost origins for local development.
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin != "" && isAllowedOrigin(origin) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		// Handle preflight
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func isAllowedOrigin(origin string) bool {
	raw := os.Getenv("ALLOWED_ORIGINS")
	if raw != "" {
		for _, o := range strings.Split(raw, ",") {
			if strings.TrimSpace(o) == origin {
				return true
			}
		}
	}
	// Default: allow localhost variants
	localOrigins := []string{
		"http://localhost:3000",
		"http://localhost:8080",
		"http://127.0.0.1:3000",
		"http://127.0.0.1:8080",
	}
	for _, o := range localOrigins {
		if o == origin {
			return true
		}
	}
	return false
}
