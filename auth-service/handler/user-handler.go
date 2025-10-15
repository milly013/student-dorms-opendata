package handler

import (
	"auth-service/model"
	"auth-service/service"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type UserHandler struct {
	userService *service.UserService
}

// NewUserHandler kreira novi handler
func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// RegisterHandler -> POST /register
func (h *UserHandler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var user model.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.userService.Register(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User registered successfully"})
}

func (h *UserHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, err := h.userService.Login(req.Email, req.Password)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Generisanje JWT tokena
	token, err := service.GenerateJWT(user.ID.Hex(), user.Role)
	if err != nil {
		http.Error(w, "Could not generate token", http.StatusInternalServerError)
		return
	}

	// LOGOVANJE podataka koji idu u token
	fmt.Printf("🔑 User logged in -> user_id: %s, role: %s\n", user.ID.Hex(), user.Role)

	// Vraćanje jednog JSON odgovora
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Login successful",
		"userId":  user.ID.Hex(),
		"role":    user.Role,
		"token":   token,
	})
}

// GetAllUsersHandler -> GET /users
func (h *UserHandler) GetAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	users, err := h.userService.GetAllUsers()
	if err != nil {
		http.Error(w, "Failed to fetch users", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(users)
}

// GetUserByIDHandler -> GET /users/{id}
func (h *UserHandler) GetUserByIDHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	user, err := h.userService.GetUserByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(user)
}

// DeleteUserHandler -> DELETE /users/{id}
func (h *UserHandler) DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if err := h.userService.DeleteUserByID(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "User deleted successfully"})
}
func (h *UserHandler) GetUserRoleByIDHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	role, err := h.userService.GetUserRole(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"role": role})
}

// POST /users/{id}/in-dorm – postavi studenta u dom
func (h *UserHandler) SetInDorm(w http.ResponseWriter, r *http.Request) {
	userID := mux.Vars(r)["id"]

	if err := h.userService.SetInDorm(userID); err != nil {
		http.Error(w, "failed to set student in dorm", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "student postavljen u dom"})
}

// POST /users/{id}/out-dorm – izbaci studenta iz doma
func (h *UserHandler) SetOutDorm(w http.ResponseWriter, r *http.Request) {
	userID := mux.Vars(r)["id"]

	if err := h.userService.SetOutDorm(userID); err != nil {
		http.Error(w, "failed to set student out of dorm", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "student izbačen iz doma"})
}

// POST /users/{id}/assign-dorm – dodaj korisnika u dom
func (h *UserHandler) AssignDorm(w http.ResponseWriter, r *http.Request) {
	userID := mux.Vars(r)["id"]

	var req struct {
		DormID string `json:"dorm_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.DormID == "" {
		http.Error(w, "dorm_id is required", http.StatusBadRequest)
		return
	}

	if err := h.userService.AssignDorm(userID, req.DormID); err != nil {
		http.Error(w, "failed to assign dorm", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "student uspešno dodeljen domu",
		"userId":  userID,
		"dormId":  req.DormID,
	})
}

// POST /users/{id}/remove-dorm – izbaci korisnika iz doma
func (h *UserHandler) RemoveDorm(w http.ResponseWriter, r *http.Request) {
	userID := mux.Vars(r)["id"]

	if err := h.userService.RemoveDorm(userID); err != nil {
		http.Error(w, "failed to remove dorm", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "student uspešno iseljen iz doma",
		"userId":  userID,
	})
}
func (h *UserHandler) GetPublicUserHandler(w http.ResponseWriter, r *http.Request) {
	userId := mux.Vars(r)["id"]

	user, err := h.userService.GetUserByID(userId)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Vrati samo javne podatke (npr. username, dorm_id, inDorm)
	publicUser := struct {
		ID       string `json:"id"`
		Username string `json:"username"`
	}{
		ID:       user.ID.Hex(),
		Username: user.Username,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(publicUser)
}
