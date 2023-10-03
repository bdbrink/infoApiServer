package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"sync"
	"time"
)

type LocationResponse struct {
	IP        string  `json:"ip"`
	City      string  `json:"city"`
	Region    string  `json:"region"`
	Country   string  `json:"country"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
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

var (
	startTime   time.Time
	dataStorage map[string]interface{}
	mu          sync.Mutex
)

func init() {
	// Initialize the data storage map
	dataStorage = make(map[string]interface{})
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

func handleRoot(w http.ResponseWriter, r *http.Request) {
	// Provide a meaningful response for the root path ("/")
	fmt.Fprintf(w, "Welcome to the server! Use '/current-time-and-location' for time and location information.")
}

func handleInfo(w http.ResponseWriter, r *http.Request) {
	// Calculate the server's current runtime
	currentTime := time.Now()
	uptime := currentTime.Sub(startTime)

	// Get information about the Go runtime
	goVersion := runtime.Version()
	goOS := runtime.GOOS
	goArch := runtime.GOARCH

	// Additional server information
	serverInfo := fmt.Sprintf("Server Information:\n"+
		" - Uptime: %s\n"+
		" - Go Version: %s\n"+
		" - OS: %s\n"+
		" - Architecture: %s\n"+
		" - Programming Language: Go", uptime.String(), goVersion, goOS, goArch)

	// Set the response content type
	w.Header().Set("Content-Type", "text/plain")

	// Write the response to the client
	fmt.Fprintf(w, serverInfo)
}

func handleUptime(w http.ResponseWriter, r *http.Request) {
	// Calculate the server's current uptime
	currentTime := time.Now()
	uptime := currentTime.Sub(startTime)

	// Set the response content type
	w.Header().Set("Content-Type", "text/plain")

	// Write the uptime as a response
	fmt.Fprintf(w, "Uptime: %s", uptime.String())
}

func handleHealthCheck(w http.ResponseWriter, r *http.Request) {
	// You can add custom health check logic here
	// For simplicity, we'll just respond with a 200 OK
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "OK")
}

func handleEndpoints(w http.ResponseWriter, r *http.Request) {
	// Define a list of available endpoints
	endpoints := []string{
		"/info - Get server information",
		"/uptime - Get server uptime",
		"/healthcheck - Perform a server health check",
		"/endpoints - List available endpoints",
	}

	// Set the response content type
	w.Header().Set("Content-Type", "text/plain")

	// Write the list of available endpoints as a response
	fmt.Fprintln(w, "Available Endpoints:")
	for _, endpoint := range endpoints {
		fmt.Fprintln(w, endpoint)
	}
}

func handleData(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// Retrieve data from the storage
		mu.Lock()
		defer mu.Unlock()

		key := r.URL.Query().Get("key")
		if val, ok := dataStorage[key]; ok {
			response, _ := json.Marshal(map[string]interface{}{"key": key, "value": val})
			w.Header().Set("Content-Type", "application/json")
			w.Write(response)
		} else {
			http.Error(w, "Key not found", http.StatusNotFound)
		}

	case http.MethodPost:
		// Store data in the storage
		mu.Lock()
		defer mu.Unlock()

		var data map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			http.Error(w, "Invalid JSON data", http.StatusBadRequest)
			return
		}

		key, ok := data["key"].(string)
		if !ok {
			http.Error(w, "Invalid key", http.StatusBadRequest)
			return
		}

		value := data["value"]
		dataStorage[key] = value

		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, "Data stored successfully")

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func main() {
	// Define a route for getting the current UTC time and location
	http.HandleFunc("/current-time-and-location", getCurrentTimeAndLocation)

	// Define a route for the root path ("/")
	http.HandleFunc("/", handleRoot)
	http.HandleFunc("/info", handleInfo)
	http.HandleFunc("/uptime", handleUptime)
	http.HandleFunc("/healthcheck", handleHealthCheck)
	http.HandleFunc("/endpoints", handleEndpoints)
	http.HandleFunc("/data", handleData)

	// Start the server
	port := 8080
	fmt.Printf("Server is running on port %d...\n", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		fmt.Println(err)
	}
}
