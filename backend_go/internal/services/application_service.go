package services

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"

	"github.com/rojgarsetu/backend/internal/db"
)

type ApplicationService struct {
	db *db.PostgresDB
}

func NewApplicationService(d *db.PostgresDB) *ApplicationService {
	return &ApplicationService{db: d}
}

func (s *ApplicationService) Apply(ctx context.Context, jobID, candidateID string, req db.ApplyJobRequest) (*db.JobApplication, error) {
	jid, err := uuid.Parse(jobID)
	if err != nil {
		return nil, fmt.Errorf("invalid job id: %w", err)
	}
	cid, err := uuid.Parse(candidateID)
	if err != nil {
		return nil, fmt.Errorf("invalid candidate id: %w", err)
	}
	result, err := s.db.Queries.CreateJobApplication(ctx, db.CreateJobApplicationParams{
		JobID:       jid,
		CandidateID: cid,
		CoverLetter: req.CoverLetter,
		ResumeUrl:   req.ResumeURL,
	})
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *ApplicationService) GetApplicationByID(ctx context.Context, id string) (*db.JobApplication, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	result, err := s.db.Queries.GetJobApplicationByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *ApplicationService) GetApplicationsByJobID(ctx context.Context, jobID string, page, limit int) ([]db.JobApplication, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit
	uid, err := uuid.Parse(jobID)
	if err != nil {
		return nil, 0, fmt.Errorf("invalid job id: %w", err)
	}
	rows, err := s.db.Queries.ListJobApplicationsByJobID(ctx, db.ListJobApplicationsByJobIDParams{
		JobID:  uid,
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, 0, err
	}
	return rows, len(rows), nil
}

func (s *ApplicationService) GetApplicationsByCandidateID(ctx context.Context, candidateID string, page, limit int) ([]db.JobApplication, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit
	uid, err := uuid.Parse(candidateID)
	if err != nil {
		return nil, 0, fmt.Errorf("invalid candidate id: %w", err)
	}
	rows, err := s.db.Queries.ListJobApplicationsByCandidateID(ctx, db.ListJobApplicationsByCandidateIDParams{
		CandidateID: uid,
		Limit:       int32(limit),
		Offset:      int32(offset),
	})
	if err != nil {
		return nil, 0, err
	}
	return rows, len(rows), nil
}

func (s *ApplicationService) UpdateApplicationStatus(ctx context.Context, id string, status, notes sql.NullString) (*db.JobApplication, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	result, err := s.db.Queries.UpdateJobApplicationStatus(ctx, db.UpdateJobApplicationStatusParams{
		ID:     uid,
		Status: status,
		Notes:  notes,
	})
	if err != nil {
		return nil, err
	}
	return &result, nil
}
