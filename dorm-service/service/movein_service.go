package service

import (
	"context"
	"dorm-service/model"
	"dorm-service/repo"
)

type MoveInService struct {
	repo *repo.MoveInRequestRepository
}

func NewMoveInService(r *repo.MoveInRequestRepository) *MoveInService {
	return &MoveInService{repo: r}
}

func (s *MoveInService) CreateRequest(ctx context.Context, req *model.MoveInRequest) error {
	return s.repo.Create(ctx, req)
}

func (s *MoveInService) GetAllRequests(ctx context.Context) ([]model.MoveInRequest, error) {
	return s.repo.GetAll(ctx)
}

func (s *MoveInService) ApproveRequest(ctx context.Context, id string) error {
	return s.repo.UpdateStatus(ctx, id, "approved")
}

func (s *MoveInService) RejectRequest(ctx context.Context, id string) error {
	return s.repo.UpdateStatus(ctx, id, "rejected")
}

func (s *MoveInService) GetStudentRequest(ctx context.Context, studentID string) (*model.MoveInRequest, error) {
	return s.repo.GetByStudent(ctx, studentID)
}
