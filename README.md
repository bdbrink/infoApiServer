# Info API Server

This repository contains a simple API server implemented in Go (Golang) that provides various informational functionalities. The server includes endpoints for retrieving current time and location, server information, uptime, health check, available endpoints, storing and retrieving data, generating random numbers, checking palindromes, calculating factorials, generating prime numbers, counting word frequencies, and checking perfect numbers.

## Endpoints and Functionalities

### 1. `/current-time-and-location`

- **Method:** GET
- **Functionality:** Provides the current UTC time and location information of the client making the request.

### 2. `/info`

- **Method:** GET
- **Functionality:** Provides information about the server, including uptime, Go version, operating system, and architecture.

### 3. `/uptime`

- **Method:** GET
- **Functionality:** Returns the server's current uptime.

### 4. `/healthcheck`

- **Method:** GET
- **Functionality:** Performs a basic health check on the server.

### 5. `/endpoints`

- **Method:** GET
- **Functionality:** Lists available endpoints in the API.

### 6. `/data`

- **Method:** POST (To store data) / GET (To retrieve data)
- **Functionality:** Allows storing and retrieving data using key-value pairs.

### 7. `/random-number`

- **Method:** POST (To generate a random number within a range) / GET (To generate a random number between 1 and 100)
- **Functionality:** Generates random numbers based on the request.

### 8. `/check-palindrome`

- **Method:** POST (To check if a word is a palindrome) / GET (To return a sample palindrome word)
- **Functionality:** Checks if a given word is a palindrome or returns a sample palindrome word.

### 9. `/factorial`

- **Method:** POST
- **Functionality:** Calculates the factorial of a given number.

### 10. `/primes`

- **Method:** POST
- **Functionality:** Generates a list of prime numbers up to a specified limit.

### 11. `/word-frequency`

- **Method:** POST
- **Functionality:** Counts the frequency of words in a given text.

### 12. `/perfect-number`

- **Method:** POST
- **Functionality:** Checks if a given number is a perfect number.

## Getting Started

1. Clone the repository.

2. Navigate to the project directory.

3. Run the Go server.

## Making Requests

You can make requests to the API endpoints using tools like `curl` or API testing tools like Postman. Make sure to use the appropriate HTTP method (GET, POST) and include any required JSON data in the request body for POST requests.

Example request using `curl`:

