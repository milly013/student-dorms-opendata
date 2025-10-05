package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"request-service/model"
	"request-service/service"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
)

type MoveInHandler struct {
	service *service.MoveInService
}

// Novi handler za request-service
func NewMoveInHandler(s *service.MoveInService) *MoveInHandler {
	return &MoveInHandler{service: s}
}

// Kreiranje novog zahtjeva (POST /requests)
func (h *MoveInHandler) CreateMoveInRequest(w http.ResponseWriter, r *http.Request) {
	var req model.MoveInRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid request body"})
		return
	}

	// Postavimo početni status i createdAt
	req.Status = "pending"
	req.CreatedAt = time.Now()
	req.RequestType = "move_in"

	ctx := r.Context()
	if err := h.service.CreateRequest(ctx, &req); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "failed to create request"})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(req)
}

// Kreiranje novog zahtjeva za iseljenje (POST /requests/move-out)
func (h *MoveInHandler) CreateMoveOutRequest(w http.ResponseWriter, r *http.Request) {
	var req model.MoveInRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid request body"})
		return
	}

	// Postavimo početne vrijednosti
	req.Status = "pending"
	req.CreatedAt = time.Now()
	req.RequestType = "move_out"

	ctx := r.Context()
	if err := h.service.CreateRequest(ctx, &req); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "failed to create move_out request"})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(req)
}

// Dohvatanje svih zahtjeva (GET /requests)
func (h *MoveInHandler) GetAllRequests(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	requests, err := h.service.GetAllRequests(ctx)
	if err != nil {
		http.Error(w, "failed to get requests", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(requests)
}

// Odobravanje zahtjeva (POST /requests/{id}/approve)
func (h *MoveInHandler) ApproveRequest(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	ctx := r.Context()

	req, err := h.service.GetRequestByID(ctx, id)
	if err != nil || req == nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "request not found"})
		return
	}

	if req.StudentID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid request, missing studentID"})
		return
	}

	// prvo update status u request-service
	if err := h.service.ApproveRequest(ctx, id); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "failed to approve request"})
		return
	}

	// odredi URL za auth-service na osnovu tipa zahtjeva
	var authServiceURL string
	switch req.RequestType {
	case "move_in":
		authServiceURL = fmt.Sprintf("http://api-gateway:8000/auth/users/%s/in-dorm", req.StudentID)
	case "move_out":
		authServiceURL = fmt.Sprintf("http://api-gateway:8000/auth/users/%s/out-dorm", req.StudentID)
	default:
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "unknown request type"})
		return
	}

	// proslijedi token auth-service-u
	token := r.Header.Get("Authorization")
	if token == "" {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "missing token"})
		return
	}
	if !strings.HasPrefix(token, "Bearer ") {
		token = "Bearer " + token
	}

	httpReq, _ := http.NewRequestWithContext(ctx, "POST", authServiceURL, nil)
	httpReq.Header.Set("Authorization", token)

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil || resp.StatusCode != http.StatusOK {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "failed to update user dorm status"})
		return
	}
	defer resp.Body.Close()

	message := "Request approved and user status updated"
	if req.RequestType == "move_out" {
		message = "Request approved and user removed from dorm"
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": message})
}

// Odbijanje zahtjeva (POST /requests/{id}/reject)
func (h *MoveInHandler) RejectRequest(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	ctx := r.Context()

	if err := h.service.RejectRequest(ctx, id); err != nil {
		http.Error(w, "failed to reject request", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// GET /requests/popular-dorms
func (h *MoveInHandler) GetPopularDorms(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	popular, err := h.service.GetPopularDorms(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "failed to get popular dorms"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(popular)
}

// Brisanje zahtjeva (DELETE /requests/{id})
func (h *MoveInHandler) DeleteRequest(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	ctx := r.Context()

	err := h.service.DeleteRequest(ctx, id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": "request not found"})
			return
		}
		if strings.Contains(err.Error(), "invalid ObjectID") {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "invalid request ID"})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "failed to delete request"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "request deleted successfully"})
}
