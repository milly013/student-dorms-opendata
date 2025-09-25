package handler

import (
	"dorm-service/model"
	"dorm-service/service"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type MoveInHandler struct {
	service *service.MoveInService
}

func NewMoveInHandler(s *service.MoveInService) *MoveInHandler {
	return &MoveInHandler{service: s}
}

// STUDENT: POST /dorms/{id}/movein-request
func (h *MoveInHandler) CreateMoveInRequest(w http.ResponseWriter, r *http.Request) {
	studentIDRaw := r.Context().Value("userID")
	roleRaw := r.Context().Value("role")

	studentID, ok1 := studentIDRaw.(string)
	role, ok2 := roleRaw.(string)

	if !ok1 || studentID == "" || !ok2 || role != "student" {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	vars := mux.Vars(r)
	dormID := vars["id"]

	req := &model.MoveInRequest{
		StudentID: studentID,
		DormID:    dormID,
	}

	if err := h.service.CreateRequest(r.Context(), req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Move-in request sent"})
}

// ADMIN: GET /dorms/movein-requests
func (h *MoveInHandler) GetAllRequests(w http.ResponseWriter, r *http.Request) {
	roleRaw := r.Context().Value("role")
	role, ok := roleRaw.(string)
	if !ok || role != "admin" {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	requests, err := h.service.GetAllRequests(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(requests)
}

// ADMIN: POST /dorms/movein-requests/{requestId}/approve
func (h *MoveInHandler) ApproveRequest(w http.ResponseWriter, r *http.Request) {
	roleRaw := r.Context().Value("role")
	role, ok := roleRaw.(string)
	if !ok || role != "admin" {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	vars := mux.Vars(r)
	reqID := vars["requestId"]

	if err := h.service.ApproveRequest(r.Context(), reqID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Request approved"})
}

// ADMIN: POST /dorms/movein-requests/{requestId}/reject
func (h *MoveInHandler) RejectRequest(w http.ResponseWriter, r *http.Request) {
	roleRaw := r.Context().Value("role")
	role, ok := roleRaw.(string)
	if !ok || role != "admin" {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	vars := mux.Vars(r)
	reqID := vars["requestId"]

	if err := h.service.RejectRequest(r.Context(), reqID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Request rejected"})
}
