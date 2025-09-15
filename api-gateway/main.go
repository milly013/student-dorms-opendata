package main

import (
	"io"
	"log"
	"net/http"
	"strings"
)

func main() {
	// Rutiranje za auth-service
	http.HandleFunc("/auth/", func(w http.ResponseWriter, r *http.Request) {
		forwardRequest(w, r, "http://auth-service:8080", "/auth")
	})

	// Rutiranje za dorm-service
	http.HandleFunc("/dorms/", func(w http.ResponseWriter, r *http.Request) {
		forwardRequest(w, r, "http://dorm-service:8081", "/dorms")
	})

	// Health check za API Gateway
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("API Gateway is running 🚀"))
	})

	log.Println("API Gateway running on port 8000...")
	log.Fatal(http.ListenAndServe(":8000", nil))
}

// Generalizovana proxy funkcija
func forwardRequest(w http.ResponseWriter, r *http.Request, targetService, prefix string) {
	// Ukloni prefix iz URL putanje
	trimmedPath := strings.TrimPrefix(r.URL.Path, prefix)

	// Ako je rezultat prazan, koristi "/"
	if trimmedPath == "" {
		trimmedPath = "/"
	}

	fullURL := targetService + trimmedPath

	// Napravi novi zahtev ka ciljnom servisu
	req, err := http.NewRequest(r.Method, fullURL, r.Body)
	if err != nil {
		http.Error(w, "Failed to create request", http.StatusInternalServerError)
		return
	}

	req.Header = r.Header

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Failed to reach service", http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()

	for k, v := range resp.Header {
		w.Header()[k] = v
	}
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}
