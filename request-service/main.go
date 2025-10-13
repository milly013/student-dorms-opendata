package main

import (
	"log"
	"net/http"
	"os"
	"request-service/db"
	"request-service/handler"
	"request-service/middleware"
	"request-service/repo"
	"request-service/service"

	"github.com/gorilla/mux"
)

func main() {
	uri := os.Getenv("MONGO_URI")
	if uri == "" {
		uri = "mongodb://mongo:27017"
	}
	dbName := "requestdb"

	client, err := db.Connect(uri, dbName)
	if err != nil {
		log.Fatalf("❌ Mongo connection error: %v", err)
	}
	log.Println("✅ Connected to MongoDB:", client.Name())

	if err := db.EnsureRequestCollection(client); err != nil {
		log.Println("⚠️ Kolekcija movein_requests možda već postoji:", err)
	} else {
		log.Println("✅ Kolekcija movein_requests kreirana sa validacijom")
	}

	requestRepo := repo.NewMoveInRequestRepository(client)

	// 🟢 Dodali authServiceURL koji servis treba
	authServiceURL := "http://auth-service:8080"
	requestService := service.NewMoveInService(requestRepo, authServiceURL)

	requestHandler := handler.NewMoveInHandler(requestService)

	r := mux.NewRouter()

	// Health check
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Request service is running 🚀"))
	}).Methods("GET")

	// Routes protected by JWT middleware
	r.Handle("/requests", middleware.JWTAuth(http.HandlerFunc(requestHandler.CreateMoveInRequest))).Methods("POST")
	r.Handle("/requests/move-out", middleware.JWTAuth(http.HandlerFunc(requestHandler.CreateMoveOutRequest))).Methods("POST")
	r.Handle("/requests", middleware.JWTAuth(http.HandlerFunc(requestHandler.GetAllRequests))).Methods("GET")
	r.Handle("/requests/{id}/approve", middleware.JWTAuth(http.HandlerFunc(requestHandler.ApproveRequest))).Methods("POST")
	r.Handle("/requests/{id}/reject", middleware.JWTAuth(http.HandlerFunc(requestHandler.RejectRequest))).Methods("POST")
	r.Handle("/requests/popular-dorms", middleware.JWTAuth(http.HandlerFunc(requestHandler.GetPopularDorms))).Methods("GET")
	r.Handle("/requests/{id}", middleware.JWTAuth(http.HandlerFunc(requestHandler.DeleteRequest))).Methods("DELETE")

	// 🧩 Nova ruta za prijavu kvara (issue report)
	r.Handle("/requests/issue", middleware.JWTAuth(http.HandlerFunc(requestHandler.CreateIssueRequestHandler))).Methods("POST")

	port := "8083"
	log.Println("Request-service running on port:", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
