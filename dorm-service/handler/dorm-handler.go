package handler

import (
	"dorm-service/model"
	"dorm-service/service"
	"encoding/json"
	"net/http"
)

type DormHandler struct {
	dormService *service.DormService
}

func NewDormHandler(dormService *service.DormService) *DormHandler {
	return &DormHandler{dormService: dormService}
}

// --- CRUD funkcije ---
func (h *DormHandler) CreateDormHandler(w http.ResponseWriter, r *http.Request) {
	var dorm model.Dorm
	if err := json.NewDecoder(r.Body).Decode(&dorm); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if err := h.dormService.CreateDorm(&dorm); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Dorm created successfully"})
}

func (h *DormHandler) GetDormHandler(w http.ResponseWriter, r *http.Request, id string) {
	dorm, err := h.dormService.GetDormByID(id)
	if err != nil || dorm == nil {
		http.Error(w, "Dorm not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(dorm)
}

func (h *DormHandler) GetAllDormsHandler(w http.ResponseWriter, r *http.Request) {
	dorms, err := h.dormService.GetAllDorms()
	if err != nil {
		http.Error(w, "Failed to fetch dorms", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(dorms)
}

func (h *DormHandler) UpdateDormHandler(w http.ResponseWriter, r *http.Request, id string) {
	var dorm model.Dorm
	if err := json.NewDecoder(r.Body).Decode(&dorm); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if err := h.dormService.UpdateDorm(id, &dorm); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"message": "Dorm updated successfully"})
}

func (h *DormHandler) DeleteDormHandler(w http.ResponseWriter, r *http.Request, id string) {
	if err := h.dormService.DeleteDorm(id); err != nil {
		http.Error(w, "Failed to delete dorm", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"message": "Dorm deleted successfully"})
}

// --- NOVO: Pretraga po gradu ---
func (h *DormHandler) SearchDormsByCityHandler(w http.ResponseWriter, r *http.Request) {
	city := r.URL.Query().Get("city")
	if city == "" {
		http.Error(w, "city query parameter is required", http.StatusBadRequest)
		return
	}
	dorms, err := h.dormService.GetDormsByCity(city)
	if err != nil {
		http.Error(w, "Failed to fetch dorms", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(dorms)
}

// --- Filter po vrsti smeštaja ---
func (h *DormHandler) FilterDormsByTypeHandler(w http.ResponseWriter, r *http.Request) {
	dormType := r.URL.Query().Get("type")
	if dormType == "" {
		http.Error(w, "type query parameter is required", http.StatusBadRequest)
		return
	}
	dorms, err := h.dormService.GetDormsByType(dormType)
	if err != nil {
		http.Error(w, "Failed to fetch dorms", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(dorms)
}
