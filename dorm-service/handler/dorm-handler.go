package handler

import (
	"dorm-service/model"
	"dorm-service/service"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type DormHandler struct {
	dormService *service.DormService
}

func NewDormHandler(dormService *service.DormService) *DormHandler {
	return &DormHandler{dormService: dormService}
}

// POST /dorms
func (h *DormHandler) CreateDormHandler(w http.ResponseWriter, r *http.Request) {
	userIDRaw := r.Context().Value("userID")
	roleRaw := r.Context().Value("role")

	userID, ok1 := userIDRaw.(string)
	role, ok2 := roleRaw.(string)

	if !ok1 || userID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	if !ok2 || role == "" {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	fmt.Printf("UserID: %s, Role: %s\n", userID, role)

	if role != "admin" {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	var dorm model.Dorm
	if err := json.NewDecoder(r.Body).Decode(&dorm); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	fmt.Printf("✅ User %s (role: %s) kreira novi dom: %+v\n", userID, role, dorm)

	if err := h.dormService.CreateDorm(&dorm); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Dorm created successfully"})
}

// GET /dorms/{id}
func (h *DormHandler) GetDormHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

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
func (h *DormHandler) UpdateDormHandler(w http.ResponseWriter, r *http.Request) {
	userIDRaw := r.Context().Value("userID")
	roleRaw := r.Context().Value("role")

	userID, ok1 := userIDRaw.(string)
	role, ok2 := roleRaw.(string)

	if !ok1 || userID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	if !ok2 || role == "" {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	if role != "admin" {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]

	var dorm model.Dorm
	if err := json.NewDecoder(r.Body).Decode(&dorm); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	fmt.Printf("✅ User %s (role: %s) update dorm: %s\n", userID, role, id)

	if err := h.dormService.UpdateDorm(id, &dorm); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Dorm updated successfully"})
}

// DELETE /dorms/{id}
func (h *DormHandler) DeleteDormHandler(w http.ResponseWriter, r *http.Request) {
	userIDRaw := r.Context().Value("userID")
	roleRaw := r.Context().Value("role")

	userID, ok1 := userIDRaw.(string)
	role, ok2 := roleRaw.(string)

	if !ok1 || userID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	if !ok2 || role == "" {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	if role != "admin" {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]

	fmt.Printf("✅ User %s (role: %s) delete dorm: %s\n", userID, role, id)

	if err := h.dormService.DeleteDorm(id); err != nil {
		http.Error(w, "Failed to delete dorm", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Dorm deleted successfully"})
}
