package crawler

import "time"

type ScrapedJob struct {
	ExternalJobID           string     `json:"external_job_id"`
	Title                   string     `json:"title"`
	Organization            string     `json:"organization"`
	JobType                 string     `json:"job_type"` // GOVT or PRIVATE
	QualificationRequired   string     `json:"qualification_required"`
	TotalVacancies          int        `json:"total_vacancies"`
	SalaryRange             string     `json:"salary_range"`
	JobLocation             string     `json:"job_location"`
	OfficialNotificationURL string     `json:"official_notification_url"`
	ApplyURL                string     `json:"apply_url"`
	PublishedAt             *time.Time `json:"published_at"`
	ApplicationDeadline     *time.Time `json:"application_deadline"`
	CategoryName            string     `json:"category_name"`
	TradeName               string     `json:"trade_name"`
	RawPayload              []byte     `json:"raw_payload"`
}

type CrawlResult struct {
	SourceName    string `json:"source_name"`
	Status        string `json:"status"` // SUCCESS, FAILED
	JobsFound     int    `json:"jobs_found"`
	JobsAdded     int    `json:"jobs_added"`
	Duplicates    int    `json:"duplicates"`
	ErrorsCount   int    `json:"errors_count"`
	ErrorMessage  string `json:"error_message,omitempty"`
	ExecutionTime int64  `json:"execution_time_ms"`
}
