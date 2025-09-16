package service

import (
	"dorm-service/model"
	"dorm-service/repo"
	"errors"
)

type DormService struct {
	repo *repo.DormRepository
}

func NewDormService(repo *repo.DormRepository) *DormService {
	return &DormService{repo: repo}
}

// --- CRUD funkcije ---

func (s *DormService) CreateDorm(dorm *model.Dorm) error {
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
