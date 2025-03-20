package main

import (
	"encoding/json"
	"log"
	"net/http"
)

// Message struct defines the JSON structure for responses
type Message struct {
	Status string `json:"status"` // Status of the request (successful/failed)
	Body   string `json:"body"`   // Body message
}

// endPointHandler handles the /ping endpoint
func endPointHandler(writer http.ResponseWriter, request *http.Request) {
	// Set response header to application/json
	writer.Header().Set("Content-Type", "application/json")

	// Set HTTP status code to 200 OK
	writer.WriteHeader(http.StatusOK)

	// Create success message
	message := Message{
		Status: "successful",
		Body:   "Hi, how can I help you?",
	}

	// Encode and send the message as JSON
	err := json.NewEncoder(writer).Encode(&message)
	if err != nil {
		// If encoding fails, send 500 Internal Server Error
		http.Error(writer, "Failed to encode message", http.StatusInternalServerError)
		log.Println("Error encoding JSON:", err)
	}
}

func main() {
	// Attach the /ping endpoint with rate limiting middleware
	// RateLimiter wraps the endPointHandler
	http.Handle("/ping", RateLimiter(endPointHandler))

	log.Println("Server started at :8080")

	// Start the HTTP server at port 8080
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Server failed:", err) // Log server failure
	}
}
