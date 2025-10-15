package handler

import (
	"bytes"
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

func (h *MoveInHandler) CreateMoveOutRequest(w http.ResponseWriter, r *http.Request) {
	var req model.MoveInRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid request body"})
		return
	}
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

func (h *MoveInHandler) ApproveRequest(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	ctx := r.Context()

	req, err := h.service.GetRequestByID(ctx, id)
	if err != nil || req == nil {
		http.Error(w, "request not found", http.StatusNotFound)
		return
	}

	if req.StudentID == "" {
		http.Error(w, "invalid request, missing studentID", http.StatusBadRequest)
		return
	}

	// Update status u bazi prvo
	if err := h.service.ApproveRequest(ctx, id); err != nil {
		http.Error(w, "failed to approve request", http.StatusInternalServerError)
		return
	}

	// Za move_in / move_out radi dodatne akcije
	token := r.Header.Get("Authorization")
	if token == "" {
		http.Error(w, "missing token", http.StatusUnauthorized)
		return
	}
	if !strings.HasPrefix(token, "Bearer ") {
		token = "Bearer " + token
	}

	client := &http.Client{}
	var authURL, dormURL string

	switch req.RequestType {
	case "move_in":
		authURL = fmt.Sprintf("http://api-gateway:8000/auth/users/%s/assign-dorm", req.StudentID)
		dormURL = fmt.Sprintf("http://api-gateway:8000/dorms/dorms/%s/add-user", req.DormID)
	case "move_out":
		authURL = fmt.Sprintf("http://api-gateway:8000/auth/users/%s/remove-dorm", req.StudentID)
		dormURL = fmt.Sprintf("http://api-gateway:8000/dorms/dorms/%s/remove-user", req.DormID)
	case "issue_report":
		// Samo update status, nema dorm/auth promjene
	default:
		// Ako je nepoznat tip, samo loguj i update status
		fmt.Println("⚠️ ApproveRequest: unknown request type, samo update status")
	}

	if req.RequestType == "move_in" || req.RequestType == "move_out" {
		bodyData := map[string]string{"dorm_id": req.DormID}
		bodyJSON, _ := json.Marshal(bodyData)

		authReq, _ := http.NewRequestWithContext(ctx, "POST", authURL, bytes.NewBuffer(bodyJSON))
		authReq.Header.Set("Authorization", token)
		authReq.Header.Set("Content-Type", "application/json")
		authResp, err := client.Do(authReq)
		if err != nil || authResp.StatusCode != http.StatusOK {
			http.Error(w, "failed to update auth user status", http.StatusInternalServerError)
			return
		}
		defer authResp.Body.Close()

		dormBody := map[string]string{"user_id": req.StudentID}
		dormBodyJSON, _ := json.Marshal(dormBody)

		dormReq, _ := http.NewRequestWithContext(ctx, "POST", dormURL, bytes.NewBuffer(dormBodyJSON))
		dormReq.Header.Set("Authorization", token)
		dormReq.Header.Set("Content-Type", "application/json")
		dormResp, err := client.Do(dormReq)
		if err != nil || dormResp.StatusCode != http.StatusOK {
			http.Error(w, "failed to update dorm user list", http.StatusInternalServerError)
			return
		}
		defer dormResp.Body.Close()
	}

	message := "Request approved successfully"
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": message})
}

func (h *MoveInHandler) RejectRequest(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	ctx := r.Context()

	if err := h.service.RejectRequest(ctx, id); err != nil {
		http.Error(w, "failed to reject request", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

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

func (h *MoveInHandler) CreateIssueRequestHandler(w http.ResponseWriter, r *http.Request) {
	var body struct {
		StudentID   string `json:"student_id"`
		DormID      string `json:"dorm_id"`
		Description string `json:"description"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Neispravan format zahtjeva", http.StatusBadRequest)
		return
	}

	token := r.Header.Get("Authorization")
	if token == "" {
		http.Error(w, "Nedostaje Authorization token", http.StatusUnauthorized)
		return
	}

	req, err := h.service.CreateIssueRequest(body.StudentID, body.DormID, body.Description, token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(req)
}

func (h *MoveInHandler) GetRequestsByTypeHandler(w http.ResponseWriter, r *http.Request) {
	requestType := mux.Vars(r)["type"]

	requests, err := h.service.GetRequestsByType(requestType)
	if err != nil {
		http.Error(w, "Failed to get requests", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(requests)
}
