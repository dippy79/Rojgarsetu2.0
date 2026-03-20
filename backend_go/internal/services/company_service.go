package services

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/rojgarsetu/backend/internal/db"
)

type CompanyService struct {
	db *db.PostgresDB
}

func NewCompanyService(d *db.PostgresDB) *CompanyService {
	return &CompanyService{db: d}
}

func (s *CompanyService) GetCompanyByUserID(ctx context.Context, userID string) (*db.Company, error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user id: %w", err)
	}
	result, err := s.db.Queries.GetCompanyByUserID(ctx, uid)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *CompanyService) GetCompanyByID(ctx context.Context, id string) (*db.Company, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	result, err := s.db.Queries.GetCompanyByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *CompanyService) UpdateCompany(ctx context.Context, id string, req db.UpdateCompanyRequest) (*db.Company, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	result, err := s.db.Queries.UpdateCompany(ctx, db.UpdateCompanyParams{
		ID:           uid,
		Name:         req.Name,
		Industry:     req.Industry,
		CompanySize:  req.CompanySize,
		Website:      req.Website,
		LogoUrl:      req.LogoUrl,
		Description:  req.Description,
		Headquarters: req.Headquarters,
		FoundedYear:  req.FoundedYear,
	})
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *CompanyService) ListCompanies(ctx context.Context, page, limit int) ([]db.Company, int, error) {
	page = clampPage(page)
	limit = clampLimit(limit)
	offset := (page - 1) * limit
	rows, err := s.db.Queries.ListCompanies(ctx, db.ListCompaniesParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, 0, err
	}
	return rows, len(rows), nil
}
