package services

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/rojgarsetu/backend/internal/db"
)

type CandidateService struct {
	db *db.PostgresDB
}

func NewCandidateService(d *db.PostgresDB) *CandidateService {
	return &CandidateService{db: d}
}

func (s *CandidateService) GetCandidateByUserID(ctx context.Context, userID string) (*db.Candidate, error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user id: %w", err)
	}
	result, err := s.db.Queries.GetCandidateByUserID(ctx, uid)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *CandidateService) GetCandidateByID(ctx context.Context, id string) (*db.Candidate, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	result, err := s.db.Queries.GetCandidateByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *CandidateService) UpdateCandidate(ctx context.Context, id string, req db.UpdateCandidateRequest) (*db.Candidate, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	result, err := s.db.Queries.UpdateCandidate(ctx, db.UpdateCandidateParams{
		ID:                uid,
		Phone:             req.Phone,
		ResumeUrl:         req.ResumeUrl,
		Skills:            req.Skills,
		ExperienceYears:   req.ExperienceYears,
		CurrentCompany:    req.CurrentCompany,
		CurrentPosition:   req.CurrentPosition,
		Location:          req.Location,
		LinkedinUrl:       req.LinkedinUrl,
		GithubUrl:         req.GithubUrl,
		PortfolioUrl:      req.PortfolioUrl,
		Bio:               req.Bio,
		IsOpenToWork:      req.IsOpenToWork,
		ExpectedSalary:    req.ExpectedSalary,
		PreferredJobType:  req.PreferredJobType,
		PreferredLocation: req.PreferredLocation,
	})
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *CandidateService) ListCandidates(ctx context.Context, location string, page, limit int) ([]db.Candidate, int, error) {
	page = clampPage(page)
	limit = clampLimit(limit)
	offset := (page - 1) * limit
	rows, err := s.db.Queries.ListCandidates(ctx, db.ListCandidatesParams{
		Column1: location,
		Limit:   int32(limit),
		Offset:  int32(offset),
	})
	if err != nil {
		return nil, 0, err
	}
	return rows, len(rows), nil
}
