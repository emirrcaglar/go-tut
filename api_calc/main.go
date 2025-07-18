/*

A stateless calculator REST api.

WINDOWS
Invoke-RestMethod -Uri "http://localhost:8080/api/{OPERATION}" -Method Post -ContentType "application/json" -Body '{"x":{x},"y":{y}}'

LINUX
curl -X POST -H "Content-Type: application/json" -d '{"x":{x},"y":{y}}' http://localhost:8080/api/{OPERATION}

*/

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type CalculatorData struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("This is the Home Handler")
	w.Write([]byte("Welcome to the home page!")) // Send response to client
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("This is the healthCheckHandler")
	w.Write([]byte("API is healthy")) // Send response to client
}

func requestParser(w http.ResponseWriter, r *http.Request) (CalculatorData, bool) {
	var data CalculatorData

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return data, false
	}

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return data, false
	}

	return data, true
}

func requestSender(w http.ResponseWriter, result float64) {
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(map[string]float64{
		"result": result,
	})
}

func addHandler(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()

	data, ok := requestParser(w, r)
	if !ok {
		return
	}

	requestSender(w, data.X+data.Y)
}

func subtractHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	data, ok := requestParser(w, r)
	if !ok {
		return
	}

	requestSender(w, data.X-data.Y)

}

func multiplyHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	data, ok := requestParser(w, r)
	if !ok {
		return
	}

	requestSender(w, data.X*data.Y)

}

func divideHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	data, ok := requestParser(w, r)
	if !ok {
		return
	}

	if data.Y == 0 {
		http.Error(w, "Division by zero not allowed", http.StatusBadRequest)
		return
	}
	requestSender(w, data.X/data.Y)

}

func main() {
	// Define routes
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/api/health", healthCheckHandler)

	http.HandleFunc("POST /api/add", addHandler)
	http.HandleFunc("POST /api/subtract", subtractHandler)
	http.HandleFunc("POST /api/multiply", multiplyHandler)
	http.HandleFunc("POST /api/divide", divideHandler)

	fmt.Println("Server running on port 8080.")
	log.Fatal(http.ListenAndServe(":8080", nil)) // Start the server
}
