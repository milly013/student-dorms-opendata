package main

import (
	"auth-service/db"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	uri := os.Getenv("MONGO_URI")
	dbName := "authdb"

	// Povezivanje sa bazom
	client, err := db.Connect(uri, dbName)
	if err != nil {
		log.Fatalf("❌ Mongo connection error: %v", err)
	}
	fmt.Println("✅ Connected to MongoDB:", client.Name())

	// Kreiranje kolekcije users sa validacijom
	err = db.EnsureUserCollection(client)
	if err != nil {
		fmt.Println("⚠️ Kolekcija users možda već postoji:", err)
	} else {
		fmt.Println("✅ Kolekcija users kreirana sa validacijom")
	}

	// 📌 Health-check endpoint
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Auth service is running 🚀"))
	})

	log.Println("Auth-service running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// func handleRegister(w http.ResponseWriter, r *http.Request) {
// 	if r.Method != http.MethodPost {
// 		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
// 		return
// 	}

// 	var u model.User
// 	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
// 		http.Error(w, "Invalid request body", http.StatusBadRequest)
// 		return
// 	}

// 	u.ID = len(users) + 1
// 	users = append(users, u)

// 	w.WriteHeader(http.StatusCreated)
// 	json.NewEncoder(w).Encode(u)
// }

// func handleLogin(w http.ResponseWriter, r *http.Request) {
// 	if r.Method != http.MethodPost {
// 		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
// 		return
// 	}

// 	var req model.User
// 	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 		http.Error(w, "Invalid request body", http.StatusBadRequest)
// 		return
// 	}

// 	for _, u := range users {
// 		if u.Username == req.Username && u.Password == req.Password {
// 			json.NewEncoder(w).Encode(map[string]string{"token": "dummy-token"})
// 			return
// 		}
// 	}

// 	http.Error(w, "Invalid credentials", http.StatusUnauthorized)
// }
