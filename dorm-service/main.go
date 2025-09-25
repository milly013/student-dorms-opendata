package main

import (
	"dorm-service/db"
	"dorm-service/handler"
	"dorm-service/middleware"
	"dorm-service/repo"
	"dorm-service/service"
	"fmt"
	"log"
	"net/http"
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

	if err := db.EnsureDormCollection(client); err != nil {
		fmt.Println("⚠️ Kolekcija dorms možda već postoji:", err)
	} else {
		fmt.Println("✅ Kolekcija dorms kreirana sa validacijom")
	}

	dormRepo := repo.NewDormRepository(client)
	dormService := service.NewDormService(dormRepo)
	dormHandler := handler.NewDormHandler(dormService)

	r := mux.NewRouter()

	// Health-check
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Dorm service is running 🚀"))
	}).Methods("GET")

	r.HandleFunc("/dorms", dormHandler.GetAllDormsHandler).Methods("GET")
	r.HandleFunc("/dorms/{id}", dormHandler.GetDormHandler).Methods("GET")

	// Ove rute su zaštićene middleware-om → samo admin može dodavati, menjati, brisati
	r.Handle("/dorms", middleware.JWTAuth(http.HandlerFunc(dormHandler.CreateDormHandler))).Methods("POST")
	r.Handle("/dorms/{id}", middleware.JWTAuth(http.HandlerFunc(dormHandler.UpdateDormHandler))).Methods("PUT")
	r.Handle("/dorms/{id}", middleware.JWTAuth(http.HandlerFunc(dormHandler.DeleteDormHandler))).Methods("DELETE")

	port := "8081"
	log.Println("Dorm-service running on port:", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
