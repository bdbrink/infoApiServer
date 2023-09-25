package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
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

type ServerInfo struct {
	ServerName    string `json:"server_name"`
	ServerIP      string `json:"server_ip"`
	CurrentTime   string `json:"current_time"`
	UserAgent     string `json:"user_agent"`
	ClientCity    string `json:"client_city"`
	ClientRegion  string `json:"client_region"`
	ClientCountry string `json:"client_country"`
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

	// Get the server information
	serverName, _ := os.Hostname()
	serverIP := "127.0.0.1" // Replace with the actual server IP

	fmt.Printf("Request from IP: %s\n", ipAddress)
	fmt.Printf("User-Agent: %s\n", userAgent)
	fmt.Printf("Location: %s, %s, %s\n", location.City, location.Region, location.Country)

	// Create a ServerInfo struct
	serverInfo := ServerInfo{
		ServerName:    serverName,
		ServerIP:      serverIP,
		CurrentTime:   timeString,
		UserAgent:     userAgent,
		ClientCity:    location.City,
		ClientRegion:  location.Region,
		ClientCountry: location.Country,
	}

	// Convert the ServerInfo struct to JSON and send it as the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(serverInfo)
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
