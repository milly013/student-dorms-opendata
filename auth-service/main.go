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
	"github.com/rs/cors"
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
	r.Handle("/role/{id}", middleware.JWTAuth(http.HandlerFunc(userHandler.GetUserRoleByIDHandler))).Methods("GET")
	r.Handle("/users/{id}/in-dorm", middleware.JWTAuth(http.HandlerFunc(userHandler.SetInDorm))).Methods("POST")
	r.Handle("/users/{id}/out-dorm", middleware.JWTAuth(http.HandlerFunc(userHandler.SetOutDorm))).Methods("POST")
	r.Handle("/users/{id}/assign-dorm", middleware.JWTAuth(http.HandlerFunc(userHandler.AssignDorm))).Methods("POST")
	r.Handle("/users/{id}/remove-dorm", middleware.JWTAuth(http.HandlerFunc(userHandler.RemoveDorm))).Methods("POST")
	r.HandleFunc("/public/users/{id}", userHandler.GetPublicUserHandler).Methods("GET")

	// Health check
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Auth service is running 🚀"))
	})

	// CORS konfiguracija
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:4200"}, // Angular front
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	})

	port := "8080"
	log.Println("Auth-service running on port:", port)
	// Wrap router sa CORS handlerom
	handler := c.Handler(r)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}
