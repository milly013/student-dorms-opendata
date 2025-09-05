package main

import (
	"auth-service/model"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

var users = []model.User{
	{ID: 1, Username: "student1", Password: "pass123"},
	{ID: 2, Username: "student2", Password: "pass456"},
}

func main() {
	// Healthcheck endpoint
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "ok")
	})

	// Endpoint za registraciju
	http.HandleFunc("/register", handleRegister)

	// Endpoint za login
	http.HandleFunc("/login", handleLogin)

	port := "8080"
	log.Printf("Auth-service listening on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func handleRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var u model.User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	u.ID = len(users) + 1
	users = append(users, u)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(u)
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req model.User
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	for _, u := range users {
		if u.Username == req.Username && u.Password == req.Password {
			json.NewEncoder(w).Encode(map[string]string{"token": "dummy-token"})
			return
		}
	}

	http.Error(w, "Invalid credentials", http.StatusUnauthorized)
}
