package main

import (
	"encoding/json"
	"net/http"

	"golang.org/x/time/rate"
)

// RateLimiter wraps an HTTP handler with rate limiting logic
// Parameters:
// - next: the actual handler function to execute if the request is allowed
func RateLimiter(next func(w http.ResponseWriter, r *http.Request)) http.Handler {
	
	// Create a new rate limiter:
	// - 2 tokens per second refill rate
	// - 4 tokens max capacity (burst size)
	limiter := rate.NewLimiter(2, 4)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Check if a token is available (Allow() tries to take a token)
		// If no token → reject request
		if !limiter.Allow() {
			// Set HTTP 429 Too Many Requests status
			w.WriteHeader(http.StatusTooManyRequests)

			// Set response content type to JSON
			w.Header().Set("Content-Type", "application/json")

			// Create the failure message
			message := Message{
				Status: "failed",
				Body:   "Request Capacity has been reached, try again later",
			}

			// Encode and send the JSON response
			err := json.NewEncoder(w).Encode(&message)
			if err != nil {
				// In case of encoding error, send Internal Server Error
				http.Error(w, "Error encoding response", http.StatusInternalServerError)
			}
			return
		}

		// If token is available, call the next handler
		next(w, r)
	})
}
