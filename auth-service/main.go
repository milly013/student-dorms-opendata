package main

import (
	"auth-service/db"
	"auth-service/handler"
	"auth-service/middleware"
	"auth-service/repo"
	"auth-service/service"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
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

	userRepo := repo.NewUserRepository(client)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	r := mux.NewRouter()

	// Rute za auth
	r.HandleFunc("/register", userHandler.RegisterHandler).Methods("POST")
	r.HandleFunc("/login", userHandler.LoginHandler).Methods("POST")

	// Rute za korisnike sa JWT middleware
	r.Handle("/users", middleware.JWTAuth(http.HandlerFunc(userHandler.GetAllUsersHandler))).Methods("GET")
	r.Handle("/users/{id}", middleware.JWTAuth(http.HandlerFunc(userHandler.GetUserByIDHandler))).Methods("GET")
	r.Handle("/users/{id}", middleware.JWTAuth(http.HandlerFunc(userHandler.DeleteUserHandler))).Methods("DELETE")

	// Health check
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Auth service is running 🚀"))
	})

	port := "8080"
	log.Println("Auth-service running on port:", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
