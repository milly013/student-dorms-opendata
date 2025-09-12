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

	// Definisanje ruta
	http.HandleFunc("/register", userHandler.RegisterHandler)
	http.HandleFunc("/login", userHandler.LoginHandler)
	http.Handle("/users", middleware.JWTAuth(http.HandlerFunc(userHandler.GetAllUsersHandler)))

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Auth service is running 🚀"))
	})

	port := "8080"
	log.Println("Auth-service running on port: ", port)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
