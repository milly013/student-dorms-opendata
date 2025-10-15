package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"request-service/model"
	"request-service/repo"
	"sort"
	"time"
)

type PopularDorm struct {
	DormID       string `json:"dorm_id"`
	RequestCount int    `json:"request_count"`
}

type MoveInService struct {
	repo           *repo.MoveInRequestRepository
	authServiceURL string
}

func NewMoveInService(r *repo.MoveInRequestRepository, authURL string) *MoveInService {
	return &MoveInService{repo: r, authServiceURL: authURL}
}

func (s *MoveInService) CreateRequest(ctx context.Context, req *model.MoveInRequest) error {
	return s.repo.Create(ctx, req)
}

func (s *MoveInService) GetAllRequests(ctx context.Context) ([]model.MoveInRequest, error) {
	return s.repo.GetAll(ctx)
}

func (s *MoveInService) GetRequestByID(ctx context.Context, id string) (*model.MoveInRequest, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *MoveInService) ApproveRequest(ctx context.Context, id string) error {
	return s.repo.UpdateStatus(ctx, id, "approved")
}

// ⭐ Nova funkcija za odbijanje sa razlogom
func (s *MoveInService) RejectRequestWithReason(ctx context.Context, id string, reason string) error {
	return s.repo.UpdateStatusWithReason(ctx, id, "rejected", reason)
}

func (s *MoveInService) RejectRequest(ctx context.Context, id string) error {
	return s.repo.UpdateStatus(ctx, id, "rejected")
}

func (s *MoveInService) GetStudentRequest(ctx context.Context, studentID string) (*model.MoveInRequest, error) {
	return s.repo.GetByStudent(ctx, studentID)
}

func (s *MoveInService) GetPopularDorms(ctx context.Context) ([]PopularDorm, error) {
	requests, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	dormCount := make(map[string]int)
	for _, r := range requests {
		if r.RequestType == "move_in" {
			dormCount[r.DormID]++
		}
	}

	result := make([]PopularDorm, 0, len(dormCount))
	for dormID, count := range dormCount {
		result = append(result, PopularDorm{DormID: dormID, RequestCount: count})
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].RequestCount > result[j].RequestCount
	})

	return result, nil
}

func (s *MoveInService) DeleteRequest(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *MoveInService) CreateIssueRequest(studentID, dormID, description, token string) (*model.MoveInRequest, error) {
	inDorm, err := s.checkStudentInDorm(studentID, token)
	if err != nil {
		return nil, err
	}
	if !inDorm {
		return nil, errors.New("student nije u domu, ne može prijaviti kvar")
	}

	req := &model.MoveInRequest{
		StudentID:   studentID,
		DormID:      dormID,
		Description: description,
		RequestType: "issue_report",
		Status:      "pending",
		CreatedAt:   time.Now(),
	}

	err = s.repo.Create(context.Background(), req)
	if err != nil {
		return nil, fmt.Errorf("neuspješno kreiranje zahtjeva: %w", err)
	}

	return req, nil
}

func (s *MoveInService) checkStudentInDorm(studentID, token string) (bool, error) {
	url := fmt.Sprintf("%s/users/%s", s.authServiceURL, studentID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return false, err
	}
	req.Header.Set("Authorization", token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false, fmt.Errorf("greska pri povezivanju sa auth-service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("auth-service vratio status: %d", resp.StatusCode)
	}

	var result struct {
		ID       string `json:"id"`
		Username string `json:"username"`
		InDorm   bool   `json:"inDorm"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return false, fmt.Errorf("greska pri parsiranju odgovora: %w", err)
	}

	return result.InDorm, nil
}

func (s *MoveInService) GetRequestsByType(requestType string) ([]model.MoveInRequest, error) {
	requests, err := s.repo.GetAll(context.Background())
	if err != nil {
		return nil, err
	}

	filtered := make([]model.MoveInRequest, 0)
	for _, r := range requests {
		if r.RequestType == requestType {
			filtered = append(filtered, r)
		}
	}

	return filtered, nil
}
