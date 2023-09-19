package main

import (
	"fmt"
	"net/http"
	"time"
)

func getCurrentTime(w http.ResponseWriter, r *http.Request) {
	// Get the current UTC time
	currentTime := time.Now().UTC()

	// Format the time as a string
	timeString := currentTime.Format("2006-01-02T15:04:05.999Z")

	// Return the time as a JSON response
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"utc_time": "%s"}`, timeString)
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
