package service

import (
	"context"
	"request-service/model"
	"request-service/repo"
)

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
