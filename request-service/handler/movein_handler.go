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
		http.Error(w, "request not found", http.StatusNotFound)
		return
	}

	if req.StudentID == "" {
		http.Error(w, "invalid request, missing studentID", http.StatusBadRequest)
		return
	}

	if err := h.service.ApproveRequest(ctx, id); err != nil {
		http.Error(w, "failed to approve request", http.StatusInternalServerError)
		return
	}

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
	default:
		http.Error(w, "unknown request type", http.StatusBadRequest)
		return
	}

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
	fmt.Println("📡 Sending request to dorm-service:", dormURL)
	fmt.Println("📦 Body:", string(dormBodyJSON))
	fmt.Println("🔐 Token:", token)
	if err != nil {
		fmt.Println("❌ Error sending dorm request:", err)
	} else {
		fmt.Println("🌐 dormResp.StatusCode =", dormResp.StatusCode)
	}
	if err != nil || dormResp.StatusCode != http.StatusOK {
		http.Error(w, "failed to update dorm user list", http.StatusInternalServerError)
		return
	}
	defer dormResp.Body.Close()
	fmt.Println("📦 Sending dormBody:", string(dormBodyJSON))

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

// 🟢 Novi endpoint: prijava kvara (issue report)
func (h *MoveInHandler) CreateIssueRequestHandler(w http.ResponseWriter, r *http.Request) {
	// Struktura tijela zahtjeva
	var body struct {
		StudentID   string `json:"student_id"`
		Description string `json:"description"`
	}

	// Dekodiraj JSON body
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Neispravan format zahtjeva", http.StatusBadRequest)
		return
	}

	// Uzmemo token iz header-a
	token := r.Header.Get("Authorization")
	if token == "" {
		http.Error(w, "Nedostaje Authorization token", http.StatusUnauthorized)
		return
	}

	// Pozovi servis sa tokenom
	req, err := h.service.CreateIssueRequest(body.StudentID, body.Description, token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Vrati rezultat kao JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(req)
}

// GET /requests/type/{type} – lista zahtjeva po tipu
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
