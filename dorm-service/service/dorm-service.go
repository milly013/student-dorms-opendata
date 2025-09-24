package service

import (
	"dorm-service/model"
	"dorm-service/repo"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type DormService struct {
	repo *repo.DormRepository
}

func NewDormService(repo *repo.DormRepository) *DormService {
	return &DormService{repo: repo}
}

func (s *DormService) CreateDorm(dorm *model.Dorm) error {

	// Možeš dodati validacije ako želiš
	if dorm.Capacity < dorm.Occupied {
		return errors.New("occupied cannot be greater than capacity")
	}
	return s.repo.Insert(dorm)
}

func (s *DormService) GetDormByID(id string) (*model.Dorm, error) {
	return s.repo.FindByID(id)
}

func (s *DormService) GetAllDorms() ([]model.Dorm, error) {
	return s.repo.FindAll()
}

func (s *DormService) UpdateDorm(id string, dorm *model.Dorm) error {
	if dorm.Capacity < dorm.Occupied {
		return errors.New("occupied cannot be greater than capacity")
	}
	return s.repo.Update(id, dorm)
}

func (s *DormService) DeleteDorm(id string) error {
	return s.repo.Delete(id)
}

func (s *DormService) CheckAdmin(userID string) (bool, error) {
	resp, err := http.Get("http://auth-service:8080/role/" + userID)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("failed to get user role")
	}

	var data map[string]string
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return false, err
	}

	role := data["role"]
	return role == "admin", nil
}
