package db

import (
	"context"
	"database/sql"
)

type FeatureRepository struct {
	DB *sql.DB
}

func NewFeatureRepository(db *sql.DB) *FeatureRepository {
	return &FeatureRepository{DB: db}
}

// InsertCompanyReview inserts a new review for a company
func (r *FeatureRepository) InsertCompanyReview(ctx context.Context, review *CompanyReview) error {
	query := `
		INSERT INTO company_reviews (company_id, candidate_id, rating, review_text)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at`

	return r.DB.QueryRowContext(
		ctx, query,
		review.CompanyID,
		review.CandidateID,
		review.Rating,
		review.ReviewText,
	).Scan(&review.ID, &review.CreatedAt)
}

// InsertJobReport registers a new report/flag for a job posting
func (r *FeatureRepository) InsertJobReport(ctx context.Context, report *JobReport) error {
	query := `
		INSERT INTO job_reports (job_id, job_type, reporter_user_id, reason, status)
		VALUES ($1, $2, $3, $4, 'PENDING')
		RETURNING id, status, created_at`

	return r.DB.QueryRowContext(
		ctx, query,
		report.JobID,
		report.JobType,
		report.ReporterUserID,
		report.Reason,
	).Scan(&report.ID, &report.Status, &report.CreatedAt)
}

// InsertCandidateInternalRating inserts or updates internal recruiter ratings
func (r *FeatureRepository) InsertCandidateInternalRating(ctx context.Context, rating *CandidateInternalRating) error {
	query := `
		INSERT INTO candidate_internal_ratings (company_id, candidate_id, private_rating, recruiter_notes)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (company_id, candidate_id) 
		DO UPDATE SET private_rating = EXCLUDED.private_rating, recruiter_notes = EXCLUDED.recruiter_notes
		RETURNING id, created_at`

	return r.DB.QueryRowContext(
		ctx, query,
		rating.CompanyID,
		rating.CandidateID,
		rating.PrivateRating,
		rating.RecruiterNotes,
	).Scan(&rating.ID, &rating.CreatedAt)
}
