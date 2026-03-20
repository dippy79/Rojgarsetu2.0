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
	ResumeUrl         sql.NullString
	Skills            []string
	ExperienceYears   sql.NullInt32
	CurrentCompany    sql.NullString
	CurrentPosition   sql.NullString
	Location          sql.NullString
	LinkedinUrl       sql.NullString
	GithubUrl         sql.NullString
	PortfolioUrl      sql.NullString
	Bio               sql.NullString
	IsOpenToWork      sql.NullBool
	ExpectedSalary    sql.NullString
	PreferredJobType  []string
	PreferredLocation []string
}

type UpdateCompanyRequest struct {
	Name         string
	Industry     sql.NullString
	CompanySize  sql.NullString
	Website      sql.NullString
	LogoUrl      sql.NullString
	Description  sql.NullString
	Headquarters sql.NullString
	FoundedYear  sql.NullInt32
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
	ApplicationUrl   sql.NullString
	ApplicationEmail sql.NullString
	IsRemote         sql.NullBool
	ExpiresAt        *string
}

type ApplyJobRequest struct {
	CoverLetter sql.NullString
	ResumeURL   sql.NullString
}
