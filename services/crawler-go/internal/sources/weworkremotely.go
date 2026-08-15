package sources

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

// WeWorkRemotelySource fetches jobs from WeWorkRemotely RSS feed.
// It consumes the public RSS feed:
//
//	https://weworkremotely.com/remote-jobs.rss
//
// which returns RSS/XML with job listings including:
//
//	{ "title", "link", "description", "pubDate" }
type WeWorkRemotelySource struct {
	BaseSource
	client *http.Client
}

// rssFeed represents the RSS feed structure
type rssFeed struct {
	XMLName xml.Name   `xml:"rss"`
	Channel rssChannel `xml:"channel"`
}

// rssChannel represents the RSS channel
type rssChannel struct {
	Title string    `xml:"title"`
	Items []rssItem `xml:"item"`
}

// rssItem represents a single RSS item
type rssItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

// NewWeWorkRemotelySource creates a WeWorkRemotely source.
func NewWeWorkRemotelySource() *WeWorkRemotelySource {
	return &WeWorkRemotelySource{
		BaseSource: BaseSource{NameStr: "weworkremotely", BaseURL: "https://weworkremotely.com"},
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Fetch retrieves remote jobs from WeWorkRemotely RSS feed.
func (s *WeWorkRemotelySource) Fetch(ctx context.Context) ([]PrivJobSource, error) {
	log.Info().Msg("Starting crawl for source: WeWorkRemotely")

	url := "https://weworkremotely.com/remote-jobs.rss"
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "RojgarSetu/2.0")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("weworkremotely RSS returned status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Parse RSS feed using standard XML parser
	var feed rssFeed
	if err := xml.Unmarshal(body, &feed); err != nil {
		return nil, fmt.Errorf("failed to parse weworkremotely RSS: %w", err)
	}

	var privJobs []PrivJobSource
	for _, item := range feed.Channel.Items {
		job := s.itemToPriv(&item)
		if isValidPrivJob(&job) {
			privJobs = append(privJobs, job)
		}
	}

	log.Info().Int("totalJobs", len(privJobs)).Msg("WeWorkRemotely fetch completed")
	return privJobs, nil
}

// itemToPriv converts a WeWorkRemotely RSS item to the shared PrivJobSource model.
func (s *WeWorkRemotelySource) itemToPriv(item *rssItem) PrivJobSource {
	job := PrivJobSource{
		Source:      "weworkremotely",
		Company:     "", // Extracted from title
		Title:       strings.TrimSpace(item.Title),
		Location:    "Remote (Global)",
		URL:         item.Link,
		JobType:     "Remote",
		Description: SanitizeString(item.Description, 500),
		CreatedAt:   time.Now(),
	}

	// Parse publication date
	if item.PubDate != "" {
		if t, err := time.Parse(time.RFC1123, item.PubDate); err == nil {
			job.PostedAt = &t
		}
	}

	// Extract company name from title (format: "Company Name - Job Title")
	if strings.Contains(item.Title, " - ") {
		parts := strings.SplitN(item.Title, " - ", 2)
		if len(parts) > 0 {
			job.Company = strings.TrimSpace(parts[0])
			job.Title = strings.TrimSpace(parts[1])
		}
	}

	// Check for India-specific remote jobs
	if strings.Contains(strings.ToLower(item.Description), "india") ||
		strings.Contains(strings.ToLower(item.Title), "india") {
		job.Location = "Remote (India)"
	}

	return job
}

// Name returns the source name.
func (s *WeWorkRemotelySource) Name() string {
	return s.NameStr
}

// FetchJobs implements the JobSource interface for compatibility with the engine
func (s *WeWorkRemotelySource) FetchJobs() ([]Job, error) {
	ctx := context.Background()
	privJobs, err := s.Fetch(ctx)
	if err != nil {
		return nil, err
	}

	var jobs []Job
	for _, privJob := range privJobs {
		job := Job{
			Title:             privJob.Title,
			CompanyOrDept:     privJob.Company,
			Location:          privJob.Location,
			QualificationReq:  privJob.JobType,
			SalaryOrPayScale:  "",
			ApplyURL:          privJob.URL,
			SourceAttribution: "Source: WeWorkRemotely RSS",
			HashChecksum:      "", // Will be set by engine
		}
		jobs = append(jobs, job)
	}

	return jobs, nil
}
