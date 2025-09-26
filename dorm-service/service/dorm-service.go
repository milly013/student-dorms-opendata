package service

import (
	"dorm-service/model"
	"dorm-service/repo"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sort"
)

type DormService struct {
	repo *repo.DormRepository
}

func NewDormService(repo *repo.DormRepository) *DormService {
	return &DormService{repo: repo}
}

// --- CRUD funkcije ---
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

// --- Pretraga po gradu ---
func (s *DormService) GetDormsByCity(city string) ([]model.Dorm, error) {
	return s.repo.FindByCity(city)
}

// --- Filter po vrsti smeštaja ---
// ✅ Ispravljeno da poziva ispravno FindByType iz repozitorijuma
func (s *DormService) GetDormsByType(dormType string) ([]model.Dorm, error) {
	if dormType == "" {
		return nil, errors.New("dormType cannot be empty")
	}
	return s.repo.FindByType(dormType)
}

// GetOccupancyStats vraća statistiku slobodnih mesta po domovima
func (s *DormService) GetOccupancyStats() ([]map[string]interface{}, error) {
	return s.repo.GetOccupancyStats()
}

// GetDormsSortedByFreeSpots vraća sve domove sortirane po broju slobodnih mesta
// order = "asc" za rastuće, "desc" za opadajuće
func (s *DormService) GetDormsSortedByFreeSpots(order string) ([]model.Dorm, error) {
	dorms, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}

	// Izračunaj slobodna mesta
	for i := range dorms {
		dorms[i].FreeSpots = dorms[i].Capacity - dorms[i].Occupied
	}

	// Sortiranje
	if order == "asc" {
		sort.Slice(dorms, func(i, j int) bool {
			return dorms[i].FreeSpots < dorms[j].FreeSpots
		})
	} else if order == "desc" {
		sort.Slice(dorms, func(i, j int) bool {
			return dorms[i].FreeSpots > dorms[j].FreeSpots
		})
	}

	return dorms, nil
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
