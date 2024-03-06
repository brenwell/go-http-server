package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/api/internal/allow", handleInternal)
	http.HandleFunc("/api/internal/restrict", handleInternal)
	http.HandleFunc("/api/namespace/allow", makeExternalRequest(os.Getenv("NGINX_A")))
	http.HandleFunc("/api/namespace/restrict", makeExternalRequest(os.Getenv("NGINX_B")))
	http.HandleFunc("/api/external/allow", makeExternalRequest("https://swapi.dev/api/people/1"))
	http.HandleFunc("/api/external/restrict", makeExternalRequest("https://jsonplaceholder.typicode.com/todos/1"))

	fmt.Println("Server is running on port 80...")
	http.ListenAndServe(":80", nil)
}

func handleInternal(w http.ResponseWriter, r *http.Request) {
	logRequest(r)
	data := map[string]interface{}{
		"source": "self",
		"message": "Hello, World!",
		"status":  http.StatusOK,
	}
	sendJSONResponse(w, data, http.StatusOK)
}

func makeExternalRequest(url string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logRequest(r)
		if url == "" {
			http.Error(w, "External service URL not provided", http.StatusInternalServerError)
			return
		}
		resp, err := http.Get(url)
		if err != nil {
			http.Error(w, "Error making request to external service", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, "Error reading response from external service", http.StatusInternalServerError)
			return
		}
		sendRawResponse(w, body, resp.StatusCode)
	}
}

func logRequest(r *http.Request) {
	log.Printf("Received request: %s %s", r.Method, r.URL.Path)
}

func sendJSONResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	response, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	sendRawResponse(w, response, statusCode)
}

func sendRawResponse(w http.ResponseWriter, data []byte, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(data)
}
