package handler

import (
	"encoding/json"
	"net/http"
	"request-service/model"
	"request-service/service"

	"github.com/gorilla/mux"
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
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	if err := h.service.CreateRequest(ctx, &req); err != nil {
		http.Error(w, "failed to create request", http.StatusInternalServerError)
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

	if err := h.service.ApproveRequest(ctx, id); err != nil {
		http.Error(w, "failed to approve request", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
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
