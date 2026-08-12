package db

import "time"

type CompanyReview struct {
    ID          string    `json:"id"`
    CompanyID   string    `json:"company_id"`
    CandidateID string    `json:"candidate_id"`
    Rating      int       `json:"rating"`
    ReviewText  string    `json:"review_text"`
    CreatedAt   time.Time `json:"created_at"`
}

type JobReport struct {
    ID             string    `json:"id"`
    JobID          string    `json:"job_id"`
    JobType        string    `json:"job_type"`
    ReporterUserID string    `json:"reporter_user_id"`
    Reason         string    `json:"reason"`
    Status         string    `json:"status"`
    CreatedAt      time.Time `json:"created_at"`
}

type CandidateInternalRating struct {
    ID             string    `json:"id"`
    CompanyID      string    `json:"company_id"`
    CandidateID    string    `json:"candidate_id"`
    PrivateRating  int       `json:"private_rating"`
    RecruiterNotes string    `json:"recruiter_notes"`
    CreatedAt      time.Time `json:"created_at"`
}
