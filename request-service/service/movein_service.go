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

// MoveInService upravlja logikom zahtjeva za useljenje
type MoveInService struct {
	repo           *repo.MoveInRequestRepository
	authServiceURL string // URL auth-servisa za provjeru InDorm
}

// Novi servis
func NewMoveInService(r *repo.MoveInRequestRepository, authURL string) *MoveInService {
	return &MoveInService{repo: r, authServiceURL: authURL}
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

// CreateIssueRequest kreira zahtjev tipa "issue_report" uz prosljeđivanje tokena
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

// Helper funkcija koja provjerava InDorm status preko auth-servisa sa JWT tokenom
func (s *MoveInService) checkStudentInDorm(studentID, token string) (bool, error) {
	url := fmt.Sprintf("%s/users/%s", s.authServiceURL, studentID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return false, err
	}

	// Dodaj Authorization header
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

// GetRequestsByType vraća sve zahtjeve određenog tipa
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
