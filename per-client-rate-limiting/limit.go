package main

import (
	"encoding/json"
	"net"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

// Client holds the rate limiter and lastSeen timestamp for each IP
type Client struct {
	limiter  *rate.Limiter // Rate limiter specific to this client
	lastSeen time.Time     // Last time this client made a request
}

// Global variables for storing client data and mutex for thread-safe access
var (
	mu      sync.Mutex
	clients = make(map[string]*Client) // Map of IP -> Client struct
)

// getLimiter returns the rate limiter for the given IP
// If it doesn't exist, it creates a new limiter and stores it
func getLimiter(ip string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	// If client already exists, update lastSeen and return existing limiter
	if client, exists := clients[ip]; exists {
		client.lastSeen = time.Now()
		return client.limiter
	}

	// New client: Create new rate limiter (1 request/sec, burst up to 3)
	limiter := rate.NewLimiter(1, 3)
	clients[ip] = &Client{limiter, time.Now()}
	return limiter
}

// cleanupClients removes clients who haven't been seen in the last 3 minutes
// This prevents the clients map from growing indefinitely
func cleanupClients() {
	for {
		time.Sleep(time.Minute) // Run cleanup every 1 minute
		mu.Lock()
		for ip, client := range clients {
			// Remove clients inactive for more than 3 minutes
			if time.Since(client.lastSeen) > 3*time.Minute {
				delete(clients, ip)
			}
		}
		mu.Unlock()
	}
}

// RateLimiter middleware applies rate limiting per client IP
// Calls next handler only if rate limit is not exceeded
func RateLimiter(nextFn func(w http.ResponseWriter, r *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract client's IP from request
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			http.Error(w, "Unable to parse IP", http.StatusInternalServerError)
			return
		}

		// Get the limiter for the IP
		limiter := getLimiter(ip)

		// Check if request is allowed under rate limit
		if !limiter.Allow() {
			w.WriteHeader(http.StatusTooManyRequests) // 429 status
			json.NewEncoder(w).Encode(Message{
				Status: "failed",
				Body:   "Rate limit exceeded. Try again later.",
			})
			return
		}

		// If allowed, proceed to next handler
		nextFn(w, r)
	})
}
