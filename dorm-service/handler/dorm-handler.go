package handler

import (
	"dorm-service/model"
	"dorm-service/service"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

type DormHandler struct {
	dormService *service.DormService
}

func NewDormHandler(dormService *service.DormService) *DormHandler {
	return &DormHandler{dormService: dormService}
}

// --- CRUD funkcije ---
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

// GET /dorms/stats - statistika slobodnih mesta po domovima
func (h *DormHandler) GetOccupancyStatsHandler(w http.ResponseWriter, r *http.Request) {
	stats, err := h.dormService.GetOccupancyStats()
	if err != nil {
		http.Error(w, "Failed to fetch occupancy stats", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(stats)
}

func (h *DormHandler) GetDormsSortedHandler(w http.ResponseWriter, r *http.Request) {
	order := r.URL.Query().Get("order")
	if order == "" {
		order = "asc" // default rastuće
	}

	dorms, err := h.dormService.GetDormsSortedByFreeSpots(order)
	if err != nil {
		http.Error(w, "Failed to get sorted dorms", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(dorms)
}

func (h *DormHandler) AddRatingHandler(w http.ResponseWriter, r *http.Request) {
	type RatingRequest struct {
		DormID string  `json:"dorm_id"`
		Score  float64 `json:"score"`
	}

	// Sigurno uzimanje userID iz konteksta
	userIDRaw := r.Context().Value("userID")
	userID, ok := userIDRaw.(string)
	if !ok || userID == "" {
		http.Error(w, "Unauthorized: userID missing", http.StatusUnauthorized)
		return
	}

	// Parsiranje body-ja
	var req RatingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.DormID == "" {
		http.Error(w, "Dorm ID is required", http.StatusBadRequest)
		return
	}

	if req.Score < 1 || req.Score > 5 {
		http.Error(w, "Score must be between 1 and 5", http.StatusBadRequest)
		return
	}

	if err := h.dormService.AddRating(req.DormID, userID, req.Score); err != nil {
		http.Error(w, fmt.Sprintf("Failed to add rating: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Rating added successfully"})
}

// POST /dorms/{id}/add-user
func (h *DormHandler) AddUserToDormHandler(w http.ResponseWriter, r *http.Request) {
	dormID := mux.Vars(r)["id"]

	if dormID == "" {
		http.Error(w, "Dorm ID is required", http.StatusBadRequest)
		return
	}

	fmt.Println("📩 AddUserToDormHandler called")
	fmt.Println("DormID:", dormID)

	// Uvijek prvo pokušaj pročitati user_id iz tijela zahtjeva
	var body struct {
		UserID string `json:"user_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		fmt.Println("❌ Error decoding body:", err)
		return
	}
	userID := strings.TrimSpace(body.UserID)

	if userID == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}
	if userID == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		fmt.Println("❌ Missing user ID in both body and token")
		return
	}

	fmt.Println("✅ Final values -> DormID:", dormID, "UserID:", userID)

	// Sad imamo i dormID i userID
	err := h.dormService.AddUserToDorm(dormID, body.UserID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to add user to dorm: %v", err), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": fmt.Sprintf("User %s added to dorm %s", body.UserID, dormID),
	})
}

// POST /dorms/{id}/remove-user
func (h *DormHandler) RemoveUserFromDormHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	dormID := vars["id"]

	userIDRaw := r.Context().Value("userID")
	userID, ok := userIDRaw.(string)
	if !ok || userID == "" {
		http.Error(w, "Unauthorized: userID missing", http.StatusUnauthorized)
		return
	}

	if dormID == "" {
		http.Error(w, "Dorm ID is required", http.StatusBadRequest)
		return
	}

	err := h.dormService.RemoveUserFromDorm(dormID, userID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to remove user from dorm: %v", err), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": fmt.Sprintf("User %s removed from dorm %s", userID, dormID),
	})
}
