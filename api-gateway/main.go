package main

import (
	"io"
	"log"
	"net/http"
)

func main() {
	// Proxy/routing funkcije
	http.HandleFunc("/auth/", func(w http.ResponseWriter, r *http.Request) {
		proxyRequest("http://auth-service:8080", w, r)
	})
	http.HandleFunc("/dorm/", func(w http.ResponseWriter, r *http.Request) {
		proxyRequest("http://dorm-service:8081", w, r)
	})
	http.HandleFunc("/opendata/", func(w http.ResponseWriter, r *http.Request) {
		proxyRequest("http://opendata-service:8082", w, r)
	})

	log.Println("API Gateway running on port 8000...")
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func proxyRequest(target string, w http.ResponseWriter, r *http.Request) {
	// ukloni /auth, /dorm, ili /opendata iz RequestURI
	path := r.URL.Path
	if len(path) > 5 && path[:5] == "/auth" {
		path = path[5:] // sada path = "/register"
	}
	req, err := http.NewRequest(r.Method, target+path, r.Body)
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
	_, _ = io.Copy(w, resp.Body)
}
