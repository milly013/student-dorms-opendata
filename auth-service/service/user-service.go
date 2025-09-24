package service

import (
	"auth-service/model"
	"auth-service/repo"
	"context"
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

func GenerateJWT(userID string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // token traje 24h
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

type UserService struct {
	repo *repo.UserRepository
}

// NewUserService kreira novi service sloj
func NewUserService(repo *repo.UserRepository) *UserService {
	return &UserService{repo: repo}
}

// Register kreira novog korisnika i hash-uje lozinku
func (s *UserService) Register(user *model.User) error {
	existingUser, err := s.repo.FindByEmail(user.Email)
	if err != nil && err != mongo.ErrNoDocuments {
		return err
	}
	if existingUser != nil {
		return errors.New("user with this email already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)

	// Ako rola nije postavljena, postavi student
	if user.Role == "" {
		user.Role = "student"
	}

	return s.repo.Insert(user)
}

// Login proverava email i password
func (s *UserService) Login(email, password string) (*model.User, error) {
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	// Provera hash-ovane lozinke
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}

// GetAllUsers vraća sve korisnike
func (s *UserService) GetAllUsers() ([]model.User, error) {
	return s.repo.FindAll()
}

// ChangePassword menja lozinku
func (s *UserService) ChangePassword(email, newPassword string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	return s.repo.UpdatePassword(email, string(hashedPassword))
}

// GetUserByID vraća korisnika po ID-u
func (s *UserService) GetUserByID(id string) (*model.User, error) {
	user, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}
	return user, nil
}

// DeleteUserByID briše korisnika po ID-u
func (s *UserService) DeleteUserByID(id string) error {
	return s.repo.DeleteByID(id)
}

// HealthCheck (opciono) - jednostavan health check
func (s *UserService) HealthCheck(ctx context.Context) string {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	return "Auth service is running 🚀"
}

// Vrati ulogu korisnika po email-u
func (s *UserService) GetUserRole(ID string) (string, error) {
	user, err := s.repo.FindByID(ID)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", errors.New("user not found")
	}
	return user.Role, nil
}

// Helper funkcije
func (s *UserService) IsAdmin(ID string) (bool, error) {
	role, err := s.GetUserRole(ID)
	if err != nil {
		return false, err
	}
	return role == "admin", nil
}

func (s *UserService) IsStudent(ID string) (bool, error) {
	role, err := s.GetUserRole(ID)
	if err != nil {
		return false, err
	}
	return role == "student", nil
}
