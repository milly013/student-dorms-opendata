package main

import (
	"fmt"
	"log"
	"net/http"
	"opendata-service/db"
	"opendata-service/handler"
	"opendata-service/repo"
	"opendata-service/service"
	"os"

	"github.com/gorilla/mux"
)

func main() {
	uri := os.Getenv("MONGO_URI")
	if uri == "" {
		uri = "mongodb://mongo:27017"
	}
	dbName := "dormdb"

	client, err := db.Connect(uri, dbName)
	if err != nil {
		log.Fatalf("❌ Mongo connection error: %v", err)
	}
	fmt.Println("✅ Connected to MongoDB:", client.Name())

	dormRepo := repo.NewDormRepository(client)
	dormService := service.NewOpenDormService(dormRepo)
	dormHandler := handler.NewOpenDormHandler(dormService)

	r := mux.NewRouter()

	// Endpoint-i
	r.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Open Data Service is running 🚀"))
	}).Methods("GET")

	r.HandleFunc("/open-dorms", dormHandler.GetAllOpenDormsHandler).Methods("GET")
	r.HandleFunc("/open-dorms/{id}", dormHandler.GetDormHandler).Methods("GET")
	r.HandleFunc("/open-dorms/filter", dormHandler.FilterHandler).Methods("GET")
	// r.HandleFunc("/open-dorms/stats", dormHandler.GetStatsHandler).Methods("GET")
	// r.HandleFunc("/open-dorms/recommend", dormHandler.RecommendHandler).Methods("GET")
	// r.HandleFunc("/open-dorms/comments", dormHandler.CommentsHandler).Methods("POST")

	port := "8082"
	log.Println("Open Data Service running on port:", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
