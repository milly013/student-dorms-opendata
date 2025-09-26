package service

import (
	"opendata-service/model"
	"opendata-service/repo"
)

type OpenDormService struct {
	repo *repo.OpenDormRepository
}

func NewOpenDormService(r *repo.OpenDormRepository) *OpenDormService {
	return &OpenDormService{repo: r}
}

// Vraća sve dormove normalizovane u OpenDorm format
func (s *OpenDormService) GetAllDorms() ([]model.OpenDorm, error) {
	dorms, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}

	for i := range dorms {
		dorms[i].OccupancyRate = 0
		if dorms[i].Capacity > 0 {
			dorms[i].OccupancyRate = float64(dorms[i].Occupied) / float64(dorms[i].Capacity)
		}
		// Ovde možeš dodati logiku za Tags, AverageRating, CommentsCount
	}

	return dorms, nil
}

// Vraća jedan dorm po ID-u, normalizovan
func (s *OpenDormService) GetDormByID(id string) (*model.OpenDorm, error) {
	dorm, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if dorm.Capacity > 0 {
		dorm.OccupancyRate = float64(dorm.Occupied) / float64(dorm.Capacity)
	}

	return dorm, nil
}

// Filter metoda po kapacitetu i tipu
func (s *OpenDormService) FilterDorms(minCapacity, maxCapacity int, dormType string) ([]model.OpenDorm, error) {
	dorms, err := s.repo.FilterByCapacityOrType(minCapacity, maxCapacity, dormType)
	if err != nil {
		return nil, err
	}

	for i := range dorms {
		if dorms[i].Capacity > 0 {
			dorms[i].OccupancyRate = float64(dorms[i].Occupied) / float64(dorms[i].Capacity)
		}
	}

	return dorms, nil
}
