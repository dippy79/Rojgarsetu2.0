package db

import "database/sql"

type RegisterRequest struct {
    Name     string
    Email    string
    Password string
    Role     string
    Phone    sql.NullString
}

type UpdateCandidateRequest struct {
    Phone             sql.NullString
    ResumeURL         sql.NullString
    Skills            []string
    ExperienceYears   sql.NullInt32
    CurrentCompany    sql.NullString
    CurrentPosition   sql.NullString
    Location          sql.NullString
    LinkedinURL       sql.NullString
    GitHubURL         sql.NullString
    PortfolioURL      sql.NullString
    Bio               sql.NullString
    IsOpenToWork      sql.NullBool
    ExpectedSalary    sql.NullString
    PreferredJobType  []string
    PreferredLocation []string
}

type CandidateProfile struct {
    ID                string
    UserID            string
    Phone             sql.NullString
    ResumeURL         sql.NullString
    Skills            []string
    ExperienceYears   sql.NullInt32
    CurrentCompany    sql.NullString
    CurrentPosition   sql.NullString
    Location          sql.NullString
    LinkedinURL       sql.NullString
    GitHubURL         sql.NullString
    PortfolioURL      sql.NullString
    Bio               sql.NullString
    IsOpenToWork      sql.NullBool
    ExpectedSalary    sql.NullString
    PreferredJobType  []string
    PreferredLocation []string
}

type UpdateCompanyRequest struct {
    Name          string
    Industry      sql.NullString
    CompanySize   sql.NullString
    Website       sql.NullString
    LogoURL       sql.NullString
    Description   sql.NullString
    Headquarters  sql.NullString
    FoundedYear   sql.NullInt32
}

type CreateJobRequest struct {
    Title            string
    Description      sql.NullString
    Requirements     sql.NullString
    Responsibilities sql.NullString
    Location         sql.NullString
    JobType          sql.NullString
    ExperienceMin    sql.NullInt32
    ExperienceMax    sql.NullInt32
    SalaryMin        sql.NullInt32
    SalaryMax        sql.NullInt32
    SalaryCurrency   sql.NullString
    Skills           []string
    Benefits         []string
    ApplicationURL   sql.NullString
    ApplicationEmail sql.NullString
    IsRemote         sql.NullBool
    ExpiresAt        *string
}

type CompanyJobWithCompany struct {
    CompanyJob
    CompanyName    string         `json:"company_name"`
    CompanyLogo    sql.NullString `json:"company_logo"`
    CompanyWebsite sql.NullString `json:"company_website"`
    IsVerified     sql.NullBool   `json:"is_verified"`
}

type ApplyJobRequest struct {
    CoverLetter sql.NullString
    ResumeURL   sql.NullString
}

type SavedJob struct {
    CandidateID string
    JobID       string
}
