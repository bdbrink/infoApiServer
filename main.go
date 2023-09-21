package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type LocationResponse struct {
	IP        string  `json:"ip"`
	City      string  `json:"city"`
	Region    string  `json:"region"`
	Country   string  `json:"country"`
	Latitude  float64 `json:"loc"`
	Longitude float64 `json:"loc"`
}

func getCurrentTimeAndLocation(w http.ResponseWriter, r *http.Request) {
	// Get the IP address of the requester
	ipAddress := r.RemoteAddr

	// Get the User-Agent header from the request
	userAgent := r.Header.Get("User-Agent")

	// Get the current UTC time
	currentTime := time.Now().UTC()
	timeString := currentTime.Format("2006-01-02T15:04:05.999Z")

	// Make a request to ipinfo.io to get the location information
	locationURL := "https://ipinfo.io/" + ipAddress + "/json"
	resp, err := http.Get(locationURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Decode the response JSON into a LocationResponse struct
	var location LocationResponse
	if err := json.NewDecoder(resp.Body).Decode(&location); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Log the IP address, User-Agent, and location
	fmt.Printf("Request from IP: %s\n", ipAddress)
	fmt.Printf("User-Agent: %s\n", userAgent)
	fmt.Printf("Location: %s, %s, %s\n", location.City, location.Region, location.Country)

	// Return the time, User-Agent, and location as a JSON response
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"utc_time": "%s", "user_agent": "%s", "location": "%s, %s, %s"}`, timeString, userAgent, location.City, location.Region, location.Country)
}

func main() {
	// Define a route for getting the current UTC time and location
	http.HandleFunc("/current-time-and-location", getCurrentTimeAndLocation)

	// Start the server
	port := 8080
	fmt.Printf("Server is running on port %d...\n", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		fmt.Println(err)
	}
}
