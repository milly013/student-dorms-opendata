package service

import (
	"context"
	"request-service/model"
	"request-service/repo"
	"sort"
)

type PopularDorm struct {
	DormID       string `json:"dorm_id"`
	RequestCount int    `json:"request_count"`
}

// MoveInService upravlja logikom zahtjeva za useljenje
type MoveInService struct {
	repo *repo.MoveInRequestRepository
}

// Novi servis
func NewMoveInService(r *repo.MoveInRequestRepository) *MoveInService {
	return &MoveInService{repo: r}
}

// CreateRequest kreira novi zahtjev za useljenje
func (s *MoveInService) CreateRequest(ctx context.Context, req *model.MoveInRequest) error {
	return s.repo.Create(ctx, req)
}

// GetAllRequests vraća sve zahtjeve
func (s *MoveInService) GetAllRequests(ctx context.Context) ([]model.MoveInRequest, error) {
	return s.repo.GetAll(ctx)
}

// GetRequestByID vraća MoveInRequest po njegovom ID-ju
func (s *MoveInService) GetRequestByID(ctx context.Context, id string) (*model.MoveInRequest, error) {
	return s.repo.GetByID(ctx, id)
}

// ApproveRequest postavlja status zahtjeva na "approved"
func (s *MoveInService) ApproveRequest(ctx context.Context, id string) error {
	return s.repo.UpdateStatus(ctx, id, "approved")
}

// RejectRequest postavlja status zahtjeva na "rejected"
func (s *MoveInService) RejectRequest(ctx context.Context, id string) error {
	return s.repo.UpdateStatus(ctx, id, "rejected")
}

// GetStudentRequest vraća zahtjev po studentID-u
func (s *MoveInService) GetStudentRequest(ctx context.Context, studentID string) (*model.MoveInRequest, error) {
	return s.repo.GetByStudent(ctx, studentID)
}

func (s *MoveInService) GetPopularDorms(ctx context.Context) ([]PopularDorm, error) {
	requests, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	// Mapira domID -> broj zahtjeva, samo za move_in
	dormCount := make(map[string]int)
	for _, r := range requests {
		if r.RequestType == "move_in" { // <- filtriramo samo move_in
			dormCount[r.DormID]++
		}
	}

	// Pretvori mapu u slice za sortiranje
	result := make([]PopularDorm, 0, len(dormCount))
	for dormID, count := range dormCount {
		result = append(result, PopularDorm{DormID: dormID, RequestCount: count})
	}

	// Sortiraj po broju zahtjeva, opadajuće
	sort.Slice(result, func(i, j int) bool {
		return result[i].RequestCount > result[j].RequestCount
	})

	return result, nil
}

// DeleteRequest briše zahtjev po ID-u
func (s *MoveInService) DeleteRequest(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
