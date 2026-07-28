package services

import (
	"context"
	"errors"
	"time"

	"github.com/rojgarsetu/backend/internal/db"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	db *db.PostgresDB
}

func NewUserService(database *db.PostgresDB) *UserService {
	return &UserService{db: database}
}

// CreateUser creates a new user
func (s *UserService) CreateUser(ctx context.Context, req db.RegisterRequest) (*db.User, error) {
	// Check if email exists
	exists, err := s.db.EmailExists(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("email already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Determine role (default to candidate)
	role := req.Role
	if role == "" {
		role = "candidate"
	}
	if role != "candidate" && role != "company" {
		role = "candidate"
	}

	// Create user
	user, err := s.db.CreateUser(ctx, db.CreateUserParams{
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		Role:         role,
		Phone:        req.Phone,
	})
	if err != nil {
		return nil, err
	}

	// Create candidate or company profile based on role
	if role == "candidate" {
		_, err = s.db.CreateCandidate(ctx, user.ID)
		if err != nil {
			log.Error().Err(err).Msg("failed to create candidate profile")
		}
		// Create gamification entry
		_, err = s.db.CreateGamification(ctx, user.ID)
		if err != nil {
			log.Error().Err(err).Msg("failed to create gamification entry")
		}
	} else if role == "company" {
		_, err = s.db.CreateCompany(ctx, user.ID, req.Name)
		if err != nil {
			log.Error().Err(err).Msg("failed to create company profile")
		}
	}

	return user, nil
}

// Login authenticates a user
func (s *UserService) Login(ctx context.Context, email, password string) (*db.User, error) {
	user, err := s.db.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Update last login
	if err := s.db.UpdateLastLogin(ctx, user.ID); err != nil {
		log.Error().Err(err).Msg("failed to update last login")
	}

	return user, nil
}

// GetUserByID retrieves a user by ID
func (s *UserService) GetUserByID(ctx context.Context, id string) (*db.User, error) {
	return s.db.GetUserByID(ctx, id)
}

// GetUserByEmail retrieves a user by email
func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*db.User, error) {
	return s.db.GetUserByEmail(ctx, email)
}

// UpdateUser updates user information
func (s *UserService) UpdateUser(ctx context.Context, id string, name, phone, avatarURL string) (*db.User, error) {
	return s.db.UpdateUser(ctx, db.UpdateUserParams{
		ID:        id,
		Name:      name,
		Phone:     phone,
		AvatarURL: avatarURL,
	})
}

// VerifyUser marks a user as verified
func (s *UserService) VerifyUser(ctx context.Context, id string) error {
	return s.db.VerifyUser(ctx, id)
}

// DeactivateUser deactivates a user account
func (s *UserService) DeactivateUser(ctx context.Context, id string) error {
	return s.db.DeactivateUser(ctx, id)
}

// ListUsers retrieves users with pagination
func (s *UserService) ListUsers(ctx context.Context, page, limit int, role string) ([]db.User, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	offset := (page - 1) * limit
	return s.db.ListUsers(ctx, db.ListUsersParams{
		Limit:  limit,
		Offset: offset,
		Role:   role,
	})
}

// ============================================
// CANDIDATE SERVICE
// ============================================

type CandidateService struct {
	db *db.PostgresDB
}

func NewCandidateService(database *db.PostgresDB) *CandidateService {
	return &CandidateService{db: database}
}

// GetCandidateByUserID retrieves candidate by user ID
func (s *CandidateService) GetCandidateByUserID(ctx context.Context, userID string) (*db.Candidate, error) {
	return s.db.GetCandidateByUserID(ctx, userID)
}

// GetCandidateByID retrieves candidate by ID
func (s *CandidateService) GetCandidateByID(ctx context.Context, id string) (*db.Candidate, error) {
	return s.db.GetCandidateByID(ctx, id)
}

// UpdateCandidate updates candidate profile
func (s *CandidateService) UpdateCandidate(ctx context.Context, id string, req db.UpdateCandidateRequest) (*db.Candidate, error) {
	return s.db.UpdateCandidate(ctx, db.UpdateCandidateParams{
		ID:                 id,
		Phone:              req.Phone,
		ResumeURL:          req.ResumeURL,
		Skills:             req.Skills,
		ExperienceYears:    req.ExperienceYears,
		CurrentCompany:     req.CurrentCompany,
		CurrentPosition:    req.CurrentPosition,
		Location:           req.Location,
		LinkedinURL:        req.LinkedinURL,
		GitHubURL:          req.GitHubURL,
		PortfolioURL:       req.PortfolioURL,
		Bio:                req.Bio,
		IsOpenToWork:       req.IsOpenToWork,
		ExpectedSalary:     req.ExpectedSalary,
		PreferredJobType:   req.PreferredJobType,
		PreferredLocation:  req.PreferredLocation,
	})
}

// GetCandidateProfile retrieves detailed candidate profile
func (s *CandidateService) GetCandidateProfile(ctx context.Context, candidateID string) (*db.CandidateProfile, error) {
	return s.db.GetCandidateProfile(ctx, candidateID)
}

// UpdateCandidateProfile updates detailed candidate profile
func (s *CandidateService) UpdateCandidateProfile(ctx context.Context, candidateID string, profile db.CandidateProfile) (*db.CandidateProfile, error) {
	return s.db.UpdateCandidateProfile(ctx, profile)
}

// SearchCandidates searches candidates with filters
func (s *CandidateService) SearchCandidates(ctx context.Context, skills []string, location string, experienceMin, experienceMax int, page, limit int) ([]db.Candidate, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	offset := (page - 1) * limit
	return s.db.SearchCandidates(ctx, db.SearchCandidatesParams{
		Skills:         skills,
		Location:      location,
		ExperienceMin:  experienceMin,
		ExperienceMax: experienceMax,
		Limit:         limit,
		Offset:        offset,
	})
}

// ============================================
// COMPANY SERVICE
// ============================================

type CompanyService struct {
	db *db.PostgresDB
}

func NewCompanyService(database *db.PostgresDB) *CompanyService {
	return &CompanyService{db: database}
}

// GetCompanyByUserID retrieves company by user ID
func (s *CompanyService) GetCompanyByUserID(ctx context.Context, userID string) (*db.Company, error) {
	return s.db.GetCompanyByUserID(ctx, userID)
}

// GetCompanyByID retrieves company by ID
func (s *CompanyService) GetCompanyByID(ctx context.Context, id string) (*db.Company, error) {
	return s.db.GetCompanyByID(ctx, id)
}

// UpdateCompany updates company profile
func (s *CompanyService) UpdateCompany(ctx context.Context, id string, req db.UpdateCompanyRequest) (*db.Company, error) {
	return s.db.UpdateCompany(ctx, db.UpdateCompanyParams{
		ID:            id,
		Name:          req.Name,
		Industry:      req.Industry,
		CompanySize:   req.CompanySize,
		Website:       req.Website,
		LogoURL:       req.LogoURL,
		Description:   req.Description,
		Headquarters:  req.Headquarters,
		FoundedYear:   req.FoundedYear,
	})
}

// VerifyCompany marks a company as verified
func (s *CompanyService) VerifyCompany(ctx context.Context, id string, approvedBy string) error {
	return s.db.VerifyCompany(ctx, id, approvedBy)
}

// ApproveCompany approves a company account
func (s *CompanyService) ApproveCompany(ctx context.Context, id string, approvedBy string) error {
	return s.db.ApproveCompany(ctx, id, approvedBy)
}

// ListCompanies retrieves companies with pagination
func (s *CompanyService) ListCompanies(ctx context.Context, page, limit int, industry string, verified *bool) ([]db.Company, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	offset := (page - 1) * limit
	return s.db.ListCompanies(ctx, db.ListCompaniesParams{
		Limit:    limit,
		Offset:   offset,
		Industry: industry,
		Verified: verified,
	})
}

// ============================================
// COMPANY JOB SERVICE
// ============================================

type CompanyJobService struct {
	db *db.PostgresDB
}

func NewCompanyJobService(database *db.PostgresDB) *CompanyJobService {
	return &CompanyJobService{db: database}
}

// CreateJob creates a new job posting
func (s *CompanyJobService) CreateJob(ctx context.Context, companyID string, req db.CreateJobRequest) (*db.CompanyJob, error) {
	var expiresAt *time.Time
	if req.ExpiresAt != nil {
		t, err := time.Parse("2006-01-02", *req.ExpiresAt)
		if err == nil {
			expiresAt = &t
		}
	}

	return s.db.CreateCompanyJob(ctx, db.CreateCompanyJobParams{
		CompanyID:        companyID,
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
		ApplicationURL:   req.ApplicationURL,
		ApplicationEmail: req.ApplicationEmail,
		IsRemote:         req.IsRemote,
		ExpiresAt:        expiresAt,
	})
}

// GetJobByID retrieves job by ID
func (s *CompanyJobService) GetJobByID(ctx context.Context, id string) (*db.CompanyJobWithCompany, error) {
	return s.db.GetCompanyJobByID(ctx, id)
}

// GetJobsByCompany retrieves jobs by company ID
func (s *CompanyJobService) GetJobsByCompany(ctx context.Context, companyID string, page, limit int) ([]db.CompanyJob, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	offset := (page - 1) * limit
	return s.db.GetCompanyJobsByCompanyID(ctx, db.GetCompanyJobsByCompanyIDParams{
		CompanyID: companyID,
		Limit:     limit,
		Offset:    offset,
	})
}

// UpdateJob updates a job posting
func (s *CompanyJobService) UpdateJob(ctx context.Context, id string, req db.CreateJobRequest) (*db.CompanyJob, error) {
	var expiresAt *time.Time
	if req.ExpiresAt != nil {
		t, err := time.Parse("2006-01-02", *req.ExpiresAt)
		if err == nil {
			expiresAt = &t
		}
	}

	return s.db.UpdateCompanyJob(ctx, db.UpdateCompanyJobParams{
		ID:              id,
		Title:           req.Title,
		Description:     req.Description,
		Requirements:    req.Requirements,
		Responsibilities: req.Responsibilities,
		Location:        req.Location,
		JobType:         req.JobType,
		ExperienceMin:   req.ExperienceMin,
		ExperienceMax:   req.ExperienceMax,
		SalaryMin:       req.SalaryMin,
		SalaryMax:       req.SalaryMax,
		SalaryCurrency:  req.SalaryCurrency,
		Skills:          req.Skills,
		Benefits:        req.Benefits,
		ApplicationURL:  req.ApplicationURL,
		ApplicationEmail: req.ApplicationEmail,
		IsRemote:        req.IsRemote,
		ExpiresAt:       expiresAt,
	})
}

// DeleteJob soft deletes a job posting
func (s *CompanyJobService) DeleteJob(ctx context.Context, id string) error {
	return s.db.DeleteCompanyJob(ctx, id)
}

// ListActiveJobs retrieves all active jobs with pagination
func (s *CompanyJobService) ListActiveJobs(ctx context.Context, page, limit int, location, jobType string, skills []string) ([]db.CompanyJobWithCompany, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	offset := (page - 1) * limit
	return s.db.ListActiveCompanyJobs(ctx, db.ListActiveCompanyJobsParams{
		Limit:    limit,
		Offset:   offset,
		Location: location,
		JobType:  jobType,
		Skills:   skills,
	})
}

// IncrementJobViews increments the view count for a job
func (s *CompanyJobService) IncrementJobViews(ctx context.Context, id string) error {
	return s.db.IncrementJobViews(ctx, id)
}

// ============================================
// JOB APPLICATION SERVICE
// ============================================

type JobApplicationService struct {
	db *db.PostgresDB
}

func NewJobApplicationService(database *db.PostgresDB) *JobApplicationService {
	return &JobApplicationService{db: database}
}

// ApplyJob applies for a job
func (s *JobApplicationService) ApplyJob(ctx context.Context, jobID, candidateID string, req db.ApplyJobRequest) (*db.JobApplication, error) {
	return s.db.CreateJobApplication(ctx, db.CreateJobApplicationParams{
		JobID:       jobID,
		CandidateID: candidateID,
		CoverLetter: req.CoverLetter,
		ResumeURL:   req.ResumeURL,
	})
}

// GetApplicationByID retrieves application by ID
func (s *JobApplicationService) GetApplicationByID(ctx context.Context, id string) (*db.JobApplication, error) {
	return s.db.GetJobApplicationByID(ctx, id)
}

// GetApplicationsByJob retrieves applications for a job
func (s *JobApplicationService) GetApplicationsByJob(ctx context.Context, jobID string, page, limit int) ([]db.JobApplication, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	offset := (page - 1) * limit
	return s.db.GetJobApplicationsByJobID(ctx, db.GetJobApplicationsByJobIDParams{
		JobID:  jobID,
		Limit:  limit,
		Offset: offset,
	})
}

// GetApplicationsByCandidate retrieves applications by candidate
func (s *JobApplicationService) GetApplicationsByCandidate(ctx context.Context, candidateID string, page, limit int) ([]db.JobApplication, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	offset := (page - 1) * limit
	return s.db.GetJobApplicationsByCandidateID(ctx, db.GetJobApplicationsByCandidateIDParams{
		CandidateID: candidateID,
		Limit:       limit,
		Offset:      offset,
	})
}

// UpdateApplicationStatus updates application status
func (s *JobApplicationService) UpdateApplicationStatus(ctx context.Context, id string, status, notes string) (*db.JobApplication, error) {
	return s.db.UpdateJobApplicationStatus(ctx, db.UpdateJobApplicationStatusParams{
		ID:     id,
		Status: status,
		Notes:  notes,
	})
}

// ============================================
// SAVED JOBS SERVICE
// ============================================

type SavedJobService struct {
	db *db.PostgresDB
}

func NewSavedJobService(database *db.PostgresDB) *SavedJobService {
	return &SavedJobService{db: database}
}

// SaveJob saves a job for a candidate
func (s *SavedJobService) SaveJob(ctx context.Context, candidateID, jobID string) (*db.SavedJob, error) {
	return s.db.CreateSavedJob(ctx, candidateID, jobID)
}

// UnsaveJob removes a saved job
func (s *SavedJobService) UnsaveJob(ctx context.Context, candidateID, jobID string) error {
	return s.db.DeleteSavedJob(ctx, candidateID, jobID)
}

// GetSavedJobs retrieves saved jobs for a candidate
func (s *SavedJobService) GetSavedJobs(ctx context.Context, candidateID string, page, limit int) ([]db.CompanyJobWithCompany, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	offset := (page - 1) * limit
	return s.db.GetSavedJobs(ctx, db.GetSavedJobsParams{
		CandidateID: candidateID,
		Limit:       limit,
		Offset:      offset,
	})
}

// IsJobSaved checks if a job is saved by a candidate
func (s *SavedJobService) IsJobSaved(ctx context.Context, candidateID, jobID string) (bool, error) {
	return s.db.IsJobSaved(ctx, candidateID, jobID)
}
