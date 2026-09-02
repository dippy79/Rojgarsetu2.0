package db

import (
    "time"

    "github.com/google/uuid"
)

type CompanyReview struct {
    ID          uuid.UUID `db:"id"`
    CompanyID   uuid.UUID `db:"company_id"`
    CandidateID uuid.UUID `db:"candidate_id"`
    Rating      int       `db:"rating"`
    ReviewText  string    `db:"review_text"`
    CreatedAt   time.Time `db:"created_at"`
}

type JobReport struct {
    ID             uuid.UUID `db:"id"`
    JobID          uuid.UUID `db:"job_id"`
    JobType        string    `db:"job_type"`
    ReporterUserID uuid.UUID `db:"reporter_user_id"`
    Reason         string    `db:"reason"`
    Status         string    `db:"status"`
    CreatedAt      time.Time `db:"created_at"`
}

type CandidateInternalRating struct {
    ID             uuid.UUID `db:"id"`
    CompanyID      uuid.UUID `db:"company_id"`
    CandidateID    uuid.UUID `db:"candidate_id"`
    PrivateRating  int       `db:"private_rating"`
    RecruiterNotes string    `db:"recruiter_notes"`
    CreatedAt      time.Time `db:"created_at"`
}