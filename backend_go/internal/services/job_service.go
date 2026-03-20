package services

import (
	"context"
	"fmt"
	"time"

	"database/sql"

	"github.com/google/uuid"

	"github.com/rojgarsetu/backend/internal/db"
)

type JobService struct {
	db *db.PostgresDB
}

func NewJobService(d *db.PostgresDB) *JobService {
	return &JobService{db: d}
}

func (s *JobService) CreateJob(ctx context.Context, companyID string, req db.CreateJobRequest) (*db.CompanyJob, error) {
	cid, err := uuid.Parse(companyID)
	if err != nil {
		return nil, fmt.Errorf("invalid company id: %w", err)
	}
	var expiresAt sql.NullTime
	if req.ExpiresAt != nil {
		t, err := time.Parse("2006-01-02", *req.ExpiresAt)
		if err != nil {
			return nil, fmt.Errorf("invalid expires_at: %w", err)
		}
		expiresAt = sql.NullTime{Time: t, Valid: true}
	}
	result, err := s.db.Queries.CreateCompanyJob(ctx, db.CreateCompanyJobParams{
		CompanyID:        cid,
		Title:            req.Title,
		Description:      req.Description,
		Requirements:     req.Requirements,
		Responsibilities: req.Responsibilities,
		Location:         req.Location,
		JobType:          req.JobType,
		ExperienceMin:    req.ExperienceMin,
		ExperienceMax:    req.ExperienceMax,
		SalaryMin:        req.SalaryMin,
		SalaryMax:        req.SalaryMax,
		SalaryCurrency:   req.SalaryCurrency,
		Skills:           req.Skills,
		Benefits:         req.Benefits,
		ApplicationUrl:   req.ApplicationUrl,
		ApplicationEmail: req.ApplicationEmail,
		IsRemote:         req.IsRemote,
		ExpiresAt:        expiresAt,
	})
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *JobService) GetJobByID(ctx context.Context, id string) (*db.CompanyJob, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	result, err := s.db.Queries.GetCompanyJobByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *JobService) GetJobsByCompanyID(ctx context.Context, companyID string, page, limit int) ([]db.CompanyJob, int, error) {
	uid, err := uuid.Parse(companyID)
	if err != nil {
		return nil, 0, fmt.Errorf("invalid company id: %w", err)
	}
	page = clampPage(page)
	limit = clampLimit(limit)
	offset := (page - 1) * limit
	rows, err := s.db.Queries.GetCompanyJobsByCompanyID(ctx, db.GetCompanyJobsByCompanyIDParams{
		CompanyID: uid,
		Limit:     int32(limit),
		Offset:    int32(offset),
	})
	if err != nil {
		return nil, 0, err
	}
	return rows, len(rows), nil
}

func (s *JobService) UpdateJob(ctx context.Context, id string, req db.CreateJobRequest) (*db.CompanyJob, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	var expiresAt sql.NullTime
	if req.ExpiresAt != nil {
		t, err := time.Parse("2006-01-02", *req.ExpiresAt)
		if err != nil {
			return nil, fmt.Errorf("invalid expires_at: %w", err)
		}
		expiresAt = sql.NullTime{Time: t, Valid: true}
	}
	result, err := s.db.Queries.UpdateCompanyJob(ctx, db.UpdateCompanyJobParams{
		ID:               uid,
		Title:            req.Title,
		Description:      req.Description,
		Requirements:     req.Requirements,
		Responsibilities: req.Responsibilities,
		Location:         req.Location,
		JobType:          req.JobType,
		ExperienceMin:    req.ExperienceMin,
		ExperienceMax:    req.ExperienceMax,
		SalaryMin:        req.SalaryMin,
		SalaryMax:        req.SalaryMax,
		SalaryCurrency:   req.SalaryCurrency,
		Skills:           req.Skills,
		Benefits:         req.Benefits,
		ApplicationUrl:   req.ApplicationUrl,
		ApplicationEmail: req.ApplicationEmail,
		IsRemote:         req.IsRemote,
		ExpiresAt:        expiresAt,
	})
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *JobService) DeleteJob(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}
	return s.db.Queries.DeleteCompanyJob(ctx, uid)
}

func (s *JobService) ListActiveJobs(ctx context.Context, location, jobType string, page, limit int) ([]db.CompanyJob, int, error) {
	page = clampPage(page)
	limit = clampLimit(limit)
	offset := (page - 1) * limit
	rows, err := s.db.Queries.ListActiveCompanyJobs(ctx, db.ListActiveCompanyJobsParams{
		Column1: location,
		Column2: jobType,
		Limit:   int32(limit),
		Offset:  int32(offset),
	})
	if err != nil {
		return nil, 0, err
	}
	return rows, len(rows), nil
}

func (s *JobService) IncrementViews(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}
	return s.db.Queries.IncrementJobViews(ctx, uid)
}
