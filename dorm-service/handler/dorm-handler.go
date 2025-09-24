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

// POST /dorms
func (h *DormHandler) CreateDormHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-ID")

	isAdmin, err := h.dormService.CheckAdmin(userID)
	if err != nil || !isAdmin {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

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

// GET /dorms/{id}
func (h *DormHandler) GetDormHandler(w http.ResponseWriter, r *http.Request, id string) {
	dorm, err := h.dormService.GetDormByID(id)
	if err != nil {
		http.Error(w, "Dorm not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(dorm)
}

// GET /dorms
func (h *DormHandler) GetAllDormsHandler(w http.ResponseWriter, r *http.Request) {
	dorms, err := h.dormService.GetAllDorms()
	if err != nil {
		http.Error(w, "Failed to fetch dorms", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(dorms)
}

// PUT /dorms/{id}
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

// DELETE /dorms/{id}
func (h *DormHandler) DeleteDormHandler(w http.ResponseWriter, r *http.Request, id string) {
	if err := h.dormService.DeleteDorm(id); err != nil {
		http.Error(w, "Failed to delete dorm", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Dorm deleted successfully"})
}
