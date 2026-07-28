package db

import "context"

type RegisterRequest struct{ Email, Password, Name string }
type CreateUserParams struct{ ID, Email, Name string }
type UpdateCandidateRequest struct{ CandidateID, Phone string }
type UpdateCompanyRequest struct{ CompanyID, Name string }
type CreateJobRequest struct{ Title, Company string }
type ApplyJobRequest struct{ JobID, Candidate string }
type Candidate struct{ ID, Name, Email string }

type PostgresDB struct{}

func (d *PostgresDB) EmailExists(ctx context.Context, email string) (bool, error) { _ = ctx; _ = email; return false, nil }
func (d *PostgresDB) CreateUser(ctx context.Context, params CreateUserParams) error { _ = ctx; _ = params; return nil }
func (d *PostgresDB) CreateCandidate(ctx context.Context, c Candidate) error { _ = ctx; _ = c; return nil }
func (d *PostgresDB) UpdateCandidate(ctx context.Context, req UpdateCandidateRequest) error { _ = ctx; _ = req; return nil }
func (d *PostgresDB) UpdateCompanyJob(ctx context.Context, req UpdateCompanyRequest) error { _ = ctx; _ = req; return nil }
func (d *PostgresDB) CreateJob(ctx context.Context, req CreateJobRequest) error { _ = ctx; _ = req; return nil }
func (d *PostgresDB) ApplyJob(ctx context.Context, req ApplyJobRequest) error { _ = ctx; _ = req; return nil }
