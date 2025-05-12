package middleware

import (
	"net/http"
	"os"
)

// AuthMiddleware is an simple HTTP middleware authentication API that checks for a valid X-API-KEY header.
// If the API key does not match the one in environment variables, it rejects the request.
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("X-API-KEY")
		expectedKey := os.Getenv("API_SECRET")

		// Reject the request if the API key is missing or incorrect
		if apiKey == "" || apiKey != expectedKey {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// If valid, pass the request to the next handler
		next.ServeHTTP(w, r)
	})
}