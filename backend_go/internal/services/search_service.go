package services

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/rojgarsetu/backend/internal/db"
)

type SearchService struct {
	database *db.PostgresDB
}

func NewSearchService(database *db.PostgresDB) *SearchService {
	return &SearchService{database: database}
}

type SearchRequest struct {
	Query    string `json:"query"`
	Page     int    `json:"page"`
	Limit    int    `json:"limit"`
	Language string `json:"language"`
}

type CompanyJobSearchResult struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Location    string   `json:"location"`
	JobType     string   `json:"job_type"`
	SalaryMin   *int32   `json:"salary_min,omitempty"`
	SalaryMax   *int32   `json:"salary_max,omitempty"`
	Skills      []string `json:"skills"`
	IsRemote    bool     `json:"is_remote"`
	CompanyName string   `json:"company_name"`
	CreatedAt   string   `json:"created_at"`
	Rank        float64  `json:"rank"`
}

type GovJobSearchResult struct {
	ID          string  `json:"id"`
	Title       string  `json:"title"`
	Department  string  `json:"department"`
	Location    string  `json:"location"`
	ApplyURL    string  `json:"apply_url"`
	LastDate    string  `json:"last_date"`
	Source      string  `json:"source"`
	Eligibility string  `json:"eligibility"`
	Vacancy     *int32  `json:"vacancy_count,omitempty"`
	Salary      string  `json:"salary"`
	CreatedAt   string  `json:"created_at"`
	Rank        float64 `json:"rank"`
}

type PrivJobSearchResult struct {
	ID          string   `json:"id"`
	Company     string   `json:"company"`
	Title       string   `json:"title"`
	Location    string   `json:"location"`
	URL         string   `json:"url"`
	Salary      string   `json:"salary"`
	Experience  string   `json:"experience"`
	JobType     string   `json:"job_type"`
	Skills      []string `json:"skills"`
	Description string   `json:"description"`
	Source      string   `json:"source"`
	CreatedAt   string   `json:"created_at"`
	Rank        float64  `json:"rank"`
}

type SearchResult struct {
	CompanyJobs []CompanyJobSearchResult `json:"company_jobs"`
	GovJobs     []GovJobSearchResult     `json:"gov_jobs"`
	PrivJobs    []PrivJobSearchResult    `json:"priv_jobs"`
	Total       int                      `json:"total"`
	Page        int                      `json:"page"`
	Limit       int                      `json:"limit"`
}

func (s *SearchService) validateAndNormalizeQuery(query string) (string, string, error) {
	query = strings.TrimSpace(query)
	if query == "" {
		return "", "", fmt.Errorf("search query is required")
	}
	if len(query) < 2 {
		return "", "", fmt.Errorf("search query must be at least 2 characters")
	}
	words := strings.Fields(query)
	tsqueryParts := make([]string, 0, len(words))
	for _, w := range words {
		w = strings.ToLower(w)
		w = strings.ReplaceAll(w, "'", "''")
		tsqueryParts = append(tsqueryParts, w+":*")
	}
	tsquery := strings.Join(tsqueryParts, " & ")
	return query, tsquery, nil
}

func (s *SearchService) Search(ctx context.Context, req SearchRequest) (*SearchResult, error) {
	page := clampPage(req.Page)
	limit := clampLimit(req.Limit)
	offset := (page - 1) * limit

	rawQuery, tsquery, err := s.validateAndNormalizeQuery(req.Query)
	if err != nil {
		return nil, err
	}

	result := &SearchResult{Page: page, Limit: limit}

	companyJobs, companyCount, err := s.searchCompanyJobs(ctx, tsquery, rawQuery, req.Language, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("company jobs search failed: %w", err)
	}
	result.CompanyJobs = companyJobs
	result.Total += companyCount

	govJobs, govCount, err := s.searchGovJobs(ctx, tsquery, rawQuery, req.Language, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("gov jobs search failed: %w", err)
	}
	result.GovJobs = govJobs
	result.Total += govCount

	privJobs, privCount, err := s.searchPrivJobs(ctx, tsquery, rawQuery, req.Language, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("priv jobs search failed: %w", err)
	}
	result.PrivJobs = privJobs
	result.Total += privCount

	return result, nil
}

func (s *SearchService) searchCompanyJobs(ctx context.Context, tsquery, rawQuery, language string, limit, offset int) ([]CompanyJobSearchResult, int, error) {
	database := s.database.GetDB()
	var total int
	err := database.QueryRowContext(ctx, `
		SELECT COUNT(*)
		FROM company_jobs cj
		LEFT JOIN companies c ON cj.company_id = c.id
		WHERE cj.is_active = true
		AND ($1::text = '' OR cj.language = $1)
		AND (
			cj.search_vector @@ to_tsquery('english', $2)
			OR cj.title ILIKE '%' || $3 || '%'
			OR cj.description ILIKE '%' || $3 || '%'
			OR cj.location ILIKE '%' || $3 || '%'
			OR c.name ILIKE '%' || $3 || '%'
		)`, language, tsquery, rawQuery).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	rows, err := database.QueryContext(ctx, `
		SELECT 
			cj.id::text,
			COALESCE(cj.title, ''),
			COALESCE(cj.description, ''),
			COALESCE(cj.location, ''),
			COALESCE(cj.job_type, ''),
			cj.salary_min,
			cj.salary_max,
			COALESCE(cj.skills::text, '{}'),
			COALESCE(cj.is_remote, false),
			COALESCE(c.name, ''),
			COALESCE(cj.created_at::text, ''),
			ts_rank(cj.search_vector, to_tsquery('english', $2)) as rank
		FROM company_jobs cj
		LEFT JOIN companies c ON cj.company_id = c.id
		WHERE cj.is_active = true
		AND ($1::text = '' OR cj.language = $1)
		AND (
			cj.search_vector @@ to_tsquery('english', $2)
			OR cj.title ILIKE '%' || $3 || '%'
			OR cj.description ILIKE '%' || $3 || '%'
			OR cj.location ILIKE '%' || $3 || '%'
			OR c.name ILIKE '%' || $3 || '%'
		)
		ORDER BY rank DESC, cj.created_at DESC
		LIMIT $4 OFFSET $5
	`, language, tsquery, rawQuery, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var results []CompanyJobSearchResult
	for rows.Next() {
		var r CompanyJobSearchResult
		var skillsRaw string
		var salaryMin, salaryMax sql.NullInt32
		var isRemote sql.NullBool
		var createdAt sql.NullString

		err := rows.Scan(
			&r.ID, &r.Title, &r.Description, &r.Location, &r.JobType,
			&salaryMin, &salaryMax, &skillsRaw, &isRemote,
			&r.CompanyName, &createdAt, &r.Rank,
		)
		if err != nil {
			return nil, 0, err
		}
		if salaryMin.Valid {
			r.SalaryMin = &salaryMin.Int32
		}
		if salaryMax.Valid {
			r.SalaryMax = &salaryMax.Int32
		}
		r.IsRemote = isRemote.Valid && isRemote.Bool
		if createdAt.Valid {
			r.CreatedAt = createdAt.String
		}
		skillsRaw = strings.Trim(skillsRaw, "{}")
		if skillsRaw != "" {
			r.Skills = strings.Split(skillsRaw, ",")
			for i := range r.Skills {
				r.Skills[i] = strings.TrimSpace(r.Skills[i])
				r.Skills[i] = strings.Trim(r.Skills[i], "\"")
			}
		} else {
			r.Skills = []string{}
		}
		results = append(results, r)
	}
	return results, total, nil
}

func (s *SearchService) searchGovJobs(ctx context.Context, tsquery, rawQuery, language string, limit, offset int) ([]GovJobSearchResult, int, error) {
	db := s.database.GetDB()
	var total int
	err := db.QueryRowContext(ctx, `
		SELECT COUNT(*)
		FROM jobs_government
		WHERE is_active = true
		AND ($1::text = '' OR language = $1)
		AND (
			search_vector @@ to_tsquery('english', $2)
			OR title ILIKE '%' || $3 || '%'
			OR department ILIKE '%' || $3 || '%'
			OR location ILIKE '%' || $3 || '%'
			OR eligibility ILIKE '%' || $3 || '%'
		)
	`, language, tsquery, rawQuery).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	rows, err := db.QueryContext(ctx, `
		SELECT 
			id::text,
			COALESCE(title, ''),
			COALESCE(department, ''),
			COALESCE(location, ''),
			COALESCE(apply_url, ''),
			COALESCE(last_date::text, ''),
			COALESCE(source, ''),
			COALESCE(eligibility, ''),
			vacancy_count,
			COALESCE(salary, ''),
			COALESCE(created_at::text, ''),
			ts_rank(search_vector, to_tsquery('english', $2)) as rank
		FROM jobs_government
		WHERE is_active = true
		AND ($1::text = '' OR language = $1)
		AND (
			search_vector @@ to_tsquery('english', $2)
			OR title ILIKE '%' || $3 || '%'
			OR department ILIKE '%' || $3 || '%'
			OR location ILIKE '%' || $3 || '%'
			OR eligibility ILIKE '%' || $3 || '%'
		)
		ORDER BY rank DESC, created_at DESC
		LIMIT $4 OFFSET $5
	`, language, tsquery, rawQuery, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var results []GovJobSearchResult
	for rows.Next() {
		var r GovJobSearchResult
		var vacancyCount sql.NullInt32
		var lastDate, createdAt sql.NullString
		err := rows.Scan(
			&r.ID, &r.Title, &r.Department, &r.Location, &r.ApplyURL,
			&lastDate, &r.Source, &r.Eligibility, &vacancyCount,
			&r.Salary, &createdAt, &r.Rank,
		)
		if err != nil {
			return nil, 0, err
		}
		if vacancyCount.Valid {
			r.Vacancy = &vacancyCount.Int32
		}
		if lastDate.Valid {
			r.LastDate = lastDate.String
		}
		if createdAt.Valid {
			r.CreatedAt = createdAt.String
		}
		results = append(results, r)
	}
	return results, total, nil
}

func (s *SearchService) searchPrivJobs(ctx context.Context, tsquery, rawQuery, language string, limit, offset int) ([]PrivJobSearchResult, int, error) {
	db := s.database.GetDB()
	var total int
	err := db.QueryRowContext(ctx, `
		SELECT COUNT(*)
		FROM jobs_private
		WHERE is_active = true
		AND ($1::text = '' OR language = $1)
		AND (
			search_vector @@ to_tsquery('english', $2)
			OR title ILIKE '%' || $3 || '%'
			OR company ILIKE '%' || $3 || '%'
			OR location ILIKE '%' || $3 || '%'
			OR description ILIKE '%' || $3 || '%'
		)
	`, language, tsquery, rawQuery).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	rows, err := db.QueryContext(ctx, `
		SELECT 
			id::text,
			COALESCE(company, ''),
			COALESCE(title, ''),
			COALESCE(location, ''),
			COALESCE(url, ''),
			COALESCE(salary, ''),
			COALESCE(experience, ''),
			COALESCE(job_type, ''),
			COALESCE(skills, '{}'),
			COALESCE(description, ''),
			COALESCE(source, ''),
			COALESCE(created_at::text, ''),
			ts_rank(search_vector, to_tsquery('english', $2)) as rank
		FROM jobs_private
		WHERE is_active = true
		AND ($1::text = '' OR language = $1)
		AND (
			search_vector @@ to_tsquery('english', $2)
			OR title ILIKE '%' || $3 || '%'
			OR company ILIKE '%' || $3 || '%'
			OR location ILIKE '%' || $3 || '%'
			OR description ILIKE '%' || $3 || '%'
		)
		ORDER BY rank DESC, created_at DESC
		LIMIT $4 OFFSET $5
	`, language, tsquery, rawQuery, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var results []PrivJobSearchResult
	for rows.Next() {
		var r PrivJobSearchResult
		var skillsRaw string
		var createdAt sql.NullString
		err := rows.Scan(
			&r.ID, &r.Company, &r.Title, &r.Location, &r.URL,
			&r.Salary, &r.Experience, &r.JobType, &skillsRaw,
			&r.Description, &r.Source, &createdAt, &r.Rank,
		)
		if err != nil {
			return nil, 0, err
		}
		if createdAt.Valid {
			r.CreatedAt = createdAt.String
		}
		skillsRaw = strings.Trim(skillsRaw, "{}")
		if skillsRaw != "" {
			r.Skills = strings.Split(skillsRaw, ",")
			for i := range r.Skills {
				r.Skills[i] = strings.TrimSpace(r.Skills[i])
				r.Skills[i] = strings.Trim(r.Skills[i], "\"")
			}
		} else {
			r.Skills = []string{}
		}
		results = append(results, r)
	}
	return results, total, nil
}
