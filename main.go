package main

import (
	"fmt"
	"net/http"
	"time"
)

func getCurrentTime(w http.ResponseWriter, r *http.Request) {
	// Get the IP address of the requester
	ipAddress := r.RemoteAddr

	// Get the User-Agent header from the request
	userAgent := r.Header.Get("User-Agent")

	// Get the current UTC time
	currentTime := time.Now().UTC()

	// Format the time as a string
	timeString := currentTime.Format("2006-01-02T15:04:05.999Z")

	// Log the IP address and User-Agent
	fmt.Printf("Request from IP: %s\n", ipAddress)
	fmt.Printf("User-Agent: %s\n", userAgent)

	// Return the time and User-Agent as a JSON response
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"utc_time": "%s", "user_agent": "%s"}`, timeString, userAgent)
}

func main() {
	// Define a route for getting the current UTC time
	http.HandleFunc("/current-time", getCurrentTime)

	// Start the server
	port := 8080
	fmt.Printf("Server is running on port %d...\n", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		fmt.Println(err)
	}
}
