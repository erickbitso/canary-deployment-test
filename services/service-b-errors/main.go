package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

func main() {
	// Set up the /hello endpoint
	http.HandleFunc("/hello", helloHandler)

	// Start the HTTP server on port 8080
	fmt.Println("Server is running on port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Failed to start server: %v\n", err)
	}
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	now := time.Now() // Get the current date and time
	dateTimeStr := now.Format("2006-01-02 15:04:05")
	fmt.Println("Received a request", dateTimeStr)

	// Seed the random number generator
	rand.Seed(time.Now().UnixNano())

	// Generate a random number between 0 and 1
	if rand.Intn(2) == 0 {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, "OK", dateTimeStr)
}
