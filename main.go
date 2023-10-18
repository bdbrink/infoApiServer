package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"runtime"
	"strings"
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
	ServerName    string  `json:"server_name"`
	ServerIP      string  `json:"server_ip"`
	CurrentTime   string  `json:"current_time"`
	UserAgent     string  `json:"user_agent"`
	ClientCity    string  `json:"client_city"`
	ClientRegion  string  `json:"client_region"`
	ClientCountry string  `json:"client_country"`
	Latitude      float64 `json:"latitude"`
	Longitude     float64 `json:"longitude"`
}

var (
	startTime   time.Time
	dataStorage map[string]interface{}
	mu          sync.Mutex
)

var quotes = []string{
	"The greatest glory in living lies not in never falling, but in rising every time we fall. - Nelson Mandela",
	"Life is what happens when you're busy making other plans. - John Lennon",
	"Get busy living or get busy dying. - Stephen King",
	"You have within you right now, everything you need to deal with whatever the world can throw at you. - Brian Tracy",
	"Life is really simple, but we insist on making it complicated. - Confucius",
}

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
		Latitude:      location.Latitude,  // Assign latitude here
		Longitude:     location.Longitude, // Assign longitude here
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

func handleJsonInput(w http.ResponseWriter, r *http.Request) {
	// Check if the request method is POST
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Decode the request body
	var request struct {
		Text string `json:"text"`
	}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid JSON data", http.StatusBadRequest)
		return
	}

	// Count the number of words in the input text
	words := strings.Fields(request.Text)
	wordCount := len(words)

	// Prepare the response
	response := map[string]interface{}{
		"word_count": wordCount,
	}

	// Convert the response to JSON and send it
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func handleCalculator(w http.ResponseWriter, r *http.Request) {
	// Check if the request method is POST
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Decode the request body
	var request struct {
		Num1     float64 `json:"num1"`
		Num2     float64 `json:"num2"`
		Operator string  `json:"operator"`
	}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid JSON data", http.StatusBadRequest)
		return
	}

	// Perform the calculation based on the operator
	var result float64
	switch request.Operator {
	case "+":
		result = request.Num1 + request.Num2
	case "-":
		result = request.Num1 - request.Num2
	case "*":
		result = request.Num1 * request.Num2
	case "/":
		if request.Num2 == 0 {
			http.Error(w, "Division by zero is not allowed", http.StatusBadRequest)
			return
		}
		result = request.Num1 / request.Num2
	default:
		http.Error(w, "Invalid operator", http.StatusBadRequest)
		return
	}

	// Prepare the response
	response := map[string]interface{}{
		"result": result,
	}

	// Convert the response to JSON and send it
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func handleQuote(w http.ResponseWriter, r *http.Request) {
	// Check if the request method is GET
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Generate a random index to select a quote from the list
	rand.Seed(time.Now().UnixNano())
	randomIndex := rand.Intn(len(quotes))

	// Prepare the response
	response := map[string]interface{}{
		"quote": quotes[randomIndex],
	}

	// Convert the response to JSON and send it
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func handleRandomNumber(w http.ResponseWriter, r *http.Request) {
	// Check if the request method is GET or POST
	if r.Method == http.MethodGet || r.Method == http.MethodPost {
		// Generate a random number within the specified range
		rand.Seed(time.Now().UnixNano())

		// If it's a GET request, generate a random number between 1 and 100
		randomNumber := rand.Intn(100) + 1

		// If it's a POST request, decode the request body and generate a random number within the specified range
		if r.Method == http.MethodPost {
			var request struct {
				Min int `json:"min"`
				Max int `json:"max"`
			}

			err := json.NewDecoder(r.Body).Decode(&request)
			if err != nil {
				http.Error(w, "Invalid JSON data", http.StatusBadRequest)
				return
			}

			randomNumber = rand.Intn(request.Max-request.Min+1) + request.Min
		}

		// Prepare the response
		response := map[string]interface{}{
			"random_number": randomNumber,
		}

		// Convert the response to JSON and send it
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleCheckPalindrome(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// Handle POST request to check if the given word is a palindrome
		var request struct {
			Word string `json:"word"`
		}

		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			http.Error(w, "Invalid JSON data", http.StatusBadRequest)
			return
		}

		// Function to check if the word is a palindrome
		isPalindrome := func(word string) bool {
			runes := []rune(word)
			length := len(runes)
			for i := 0; i < length/2; i++ {
				if runes[i] != runes[length-1-i] {
					return false
				}
			}
			return true
		}

		// Check if the given word is a palindrome
		result := isPalindrome(request.Word)

		// Prepare the response
		response := map[string]interface{}{
			"is_palindrome": result,
		}

		// Convert the response to JSON and send it
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	} else if r.Method == http.MethodGet {
		// Handle GET request to return a sample palindrome word
		palindrome := "radar"

		// Prepare the response
		response := map[string]interface{}{
			"palindrome": palindrome,
		}

		// Convert the response to JSON and send it
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func main() {

	// Define a route for the root path ("/")
	http.HandleFunc("/", handleRoot)
	http.HandleFunc("/info", handleInfo)
	http.HandleFunc("/uptime", handleUptime)
	http.HandleFunc("/healthcheck", handleHealthCheck)
	http.HandleFunc("/endpoints", handleEndpoints)
	http.HandleFunc("/data", handleData)
	http.HandleFunc("/input", handleJsonInput)
	http.HandleFunc("/calculate", handleCalculator)
	http.HandleFunc("/quote", handleQuote)
	http.HandleFunc("/current-time-and-location", getCurrentTimeAndLocation)
	http.HandleFunc("/random-number", handleRandomNumber)
	http.HandleFunc("/check-palindrome", handleCheckPalindrome)

	// Start the server
	port := 8080
	fmt.Printf("Server is running on port %d...\n", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		fmt.Println(err)
	}
}
