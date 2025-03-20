package main

import (
	"encoding/json"
	"net/http"

	tollbooth "github.com/didip/tollbooth/v7"
)

// RateLimiter wraps an HTTP handler with rate limiting logic
// Parameters:
// - next: the actual handler function to execute if the request is allowed
func RateLimiter(next func(w http.ResponseWriter, r *http.Request)) http.Handler {
	// Create a limiter that allows 1 request per second
	limiter := tollbooth.NewLimiter(1, nil)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if the request is allowed
		httpError := tollbooth.LimitByRequest(limiter, w, r)
		if httpError != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(httpError.StatusCode)
			json.NewEncoder(w).Encode(map[string]string{
				"error": httpError.Message,
			})
			return
		}
		// Proceed to the actual handler
		next(w, r)
	})
}
