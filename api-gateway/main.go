package main

import (
	"io"
	"log"
	"net/http"
	"strings"
)

// Middleware za CORS
func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:4200")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Odgovor na preflight OPTIONS zahtjev
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	mux := http.NewServeMux()

	// Rutiranje za auth-service
	mux.HandleFunc("/auth/", func(w http.ResponseWriter, r *http.Request) {
		forwardRequest(w, r, "http://auth-service:8080", "/auth")
	})

	// Rutiranje za dorm-service
	mux.HandleFunc("/dorms/", func(w http.ResponseWriter, r *http.Request) {
		forwardRequest(w, r, "http://dorm-service:8081", "/dorms")
	})

	// Rutiranje za opendata-service
	mux.HandleFunc("/opendata/", func(w http.ResponseWriter, r *http.Request) {
		forwardRequest(w, r, "http://opendata-service:8082", "/opendata")
	})

	// Health check za API Gateway
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("API Gateway is running 🚀"))
	})

	log.Println("API Gateway running on port 8000...")
	log.Fatal(http.ListenAndServe(":8000", enableCORS(mux)))
}

// Generalizovana proxy funkcija
func forwardRequest(w http.ResponseWriter, r *http.Request, targetService, prefix string) {
	trimmedPath := strings.TrimPrefix(r.URL.Path, prefix)
	if trimmedPath == "" {
		trimmedPath = "/"
	}

	fullURL := targetService + trimmedPath

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
