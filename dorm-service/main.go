package main

import (
	"dorm-service/db"
	"dorm-service/handler"
	"dorm-service/repo"
	"dorm-service/service"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	uri := os.Getenv("MONGO_URI")
	if uri == "" {
		uri = "mongodb://mongo:27017" // obavezno koristi "mongo" jer je ime servisa u docker-compose-u
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

	http.HandleFunc("/dorms", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			dormHandler.GetAllDormsHandler(w, r)
		case http.MethodPost:
			dormHandler.CreateDormHandler(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/dorms/", func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Path[len("/dorms/"):]
		switch r.Method {
		case http.MethodGet:
			dormHandler.GetDormHandler(w, r, id)
		case http.MethodPut:
			dormHandler.UpdateDormHandler(w, r, id)
		case http.MethodDelete:
			dormHandler.DeleteDormHandler(w, r, id)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Dorm service is running 🚀"))
	})

	port := "8081"
	log.Println("Dorm-service running on port:", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
