package handler

import (
	"encoding/json"
	"net/http"
	"opendata-service/service"
	"strconv"

	"github.com/gorilla/mux"
)

type OpenDormHandler struct {
	service *service.OpenDormService
}

func NewOpenDormHandler(s *service.OpenDormService) *OpenDormHandler {
	return &OpenDormHandler{service: s}
}

// GET /open-dorms
func (h *OpenDormHandler) GetAllOpenDormsHandler(w http.ResponseWriter, r *http.Request) {
	dorms, err := h.service.GetAllDorms()
	if err != nil {
		http.Error(w, "Failed to fetch dorms", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dorms)
}

// GET /open-dorms/{id}
func (h *OpenDormHandler) GetDormHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	dorm, err := h.service.GetDormByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dorm)
}

// GET /open-dorms/filter?minCapacity=...&maxCapacity=...&type=...
func (h *OpenDormHandler) FilterHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	minCap, _ := strconv.Atoi(q.Get("minCapacity"))
	maxCap, _ := strconv.Atoi(q.Get("maxCapacity"))
	dType := q.Get("type")

	dorms, err := h.service.FilterDorms(minCap, maxCap, dType)
	if err != nil {
		http.Error(w, "Failed to filter dorms", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dorms)
}

// GET /healthz
func (h *OpenDormHandler) HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Open Data Service is running 🚀"))
}

func (h *OpenDormHandler) GetAveragePriceHandler(w http.ResponseWriter, r *http.Request) {
	data, err := h.service.GetAveragePricePerCity()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func (h *OpenDormHandler) GetFreeSpotsPerCityHandler(w http.ResponseWriter, r *http.Request) {
	data, err := h.service.GetFreeSpotsPerCity()
	if err != nil {
		http.Error(w, "Error fetching free spots", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func (h *OpenDormHandler) GetOccupancyPerCityHandler(w http.ResponseWriter, r *http.Request) {
	occupancy, err := h.service.GetOccupancyPerCity()
	if err != nil {
		http.Error(w, "Failed to get occupancy data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(occupancy)
}
func (h *OpenDormHandler) GetFacilitiesSummaryHandler(w http.ResponseWriter, r *http.Request) {
	summary, err := h.service.GetFacilitiesSummary()
	if err != nil {
		http.Error(w, "Greška pri dohvaćanju sadržaja", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(summary)
}
