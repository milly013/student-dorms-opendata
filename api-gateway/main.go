package main

import (
	"bytes"
	"io"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	// -----------------------
	// Auth-service rutiranje
	// -----------------------
	mux.HandleFunc("/auth/", func(w http.ResponseWriter, r *http.Request) {
		forwardRequest(w, r, "http://auth-service:8080")
	})

	// -----------------------
	// Dorm-service rutiranje (CRUD)
	// -----------------------
	mux.HandleFunc("/dorms", func(w http.ResponseWriter, r *http.Request) {
		forwardRequest(w, r, "http://dorm-service:8081")
	})
	mux.HandleFunc("/dorms/", func(w http.ResponseWriter, r *http.Request) {
		forwardRequest(w, r, "http://dorm-service:8081")
	})

	// -----------------------
	// NOVO: Filter po vrsti smeštaja
	// GET /dorms/filter?type=muški/ženski/mešoviti
	// -----------------------
	mux.HandleFunc("/dorms/filter", func(w http.ResponseWriter, r *http.Request) {
		forwardRequest(w, r, "http://dorm-service:8081")
	})

	// -----------------------
	// NOVO: Pretraga domova po gradu
	// GET /dorms/search?city=Beograd
	// -----------------------
	mux.HandleFunc("/dorms/search", func(w http.ResponseWriter, r *http.Request) {
		forwardRequest(w, r, "http://dorm-service:8081")
	})

	// -----------------------
	// Health check
	// -----------------------
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("API Gateway is running 🚀"))
	})

	log.Println("API Gateway running on port 8000...")
	log.Fatal(http.ListenAndServe(":8000", mux))
}

// -----------------------
// Generalna proxy funkcija
// -----------------------
func forwardRequest(w http.ResponseWriter, r *http.Request, targetService string) {
	// Kreiranje punog URL-a ka ciljnom servisu
	url := targetService + r.URL.Path
	if r.URL.RawQuery != "" {
		url += "?" + r.URL.RawQuery
	}

	var body []byte
	if r.Body != nil {
		body, _ = io.ReadAll(r.Body)
	}

	// Napravi novi HTTP zahtev
	req, err := http.NewRequest(r.Method, url, bytes.NewBuffer(body))
	if err != nil {
		http.Error(w, "Failed to create request", http.StatusInternalServerError)
		return
	}

	req.Header = r.Header

	// Pošalji zahtev
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Failed to reach service", http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()

	// Prosledi response
	for k, v := range resp.Header {
		w.Header()[k] = v
	}
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}
