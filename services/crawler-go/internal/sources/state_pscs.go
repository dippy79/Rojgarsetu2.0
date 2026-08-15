package sources

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/rs/zerolog/log"
)

// StatePSCSource fetches jobs from State Public Service Commissions and Subordinate Boards.
// It supports multiple states: UPPSC, UPSSSC, BPSC, MPPSC, MPSC, RPSC, RSMSSB, DSSSB, HSSC, KPSC, TNPSC, WBPSC.
type StatePSCSource struct {
	BaseSource
	httpClient *http.Client
	states     []stateConfig
}

// stateConfig holds configuration for each state commission
type stateConfig struct {
	Name        string
	URL         string
	Selector    string // CSS selector for job listings
	TitleSelect string // CSS selector for job title
	LinkSelect  string // CSS selector for job link
}

// NewStatePSCSource creates a State PSC source with multi-state support.
func NewStatePSCSource() *StatePSCSource {
	return &StatePSCSource{
		BaseSource: BaseSource{
			NameStr: "state_pscs",
			BaseURL: "",
		},
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		states: []stateConfig{
			{
				Name:        "UPPSC",
				URL:         "https://uppsc.up.nic.in/notifications.aspx",
				Selector:    ".notification-item, .notice-item, tr",
				TitleSelect: ".title, .notification-title, td:first-child",
				LinkSelect:  "a",
			},
			{
				Name:        "UPSSSC",
				URL:         "https://upsssc.gov.in/Notifications",
				Selector:    ".notification, .notice, tr",
				TitleSelect: ".title, .notice-title, td:first-child",
				LinkSelect:  "a",
			},
			{
				Name:        "BPSC",
				URL:         "https://bpsc.bih.nic.in/Advt.htm",
				Selector:    ".advertisement, .notice, tr",
				TitleSelect: ".title, .adv-title, td:first-child",
				LinkSelect:  "a",
			},
			{
				Name:        "MPPSC",
				URL:         "https://mppsc.mp.gov.in/notifications",
				Selector:    ".notification, .notice, tr",
				TitleSelect: ".title, .notice-title, td:first-child",
				LinkSelect:  "a",
			},
			{
				Name:        "MPSC",
				URL:         "https://mpsc.gov.in/advertisements",
				Selector:    ".advertisement, .notice, tr",
				TitleSelect: ".title, .adv-title, td:first-child",
				LinkSelect:  "a",
			},
			{
				Name:        "RPSC",
				URL:         "https://rpsc.rajasthan.gov.in/Advt",
				Selector:    ".advertisement, .notice, tr",
				TitleSelect: ".title, .adv-title, td:first-child",
				LinkSelect:  "a",
			},
			{
				Name:        "RSMSSB",
				URL:         "https://rsmssb.rajasthan.gov.in/notifications",
				Selector:    ".notification, .notice, tr",
				TitleSelect: ".title, .notice-title, td:first-child",
				LinkSelect:  "a",
			},
			{
				Name:        "DSSSB",
				URL:         "https://dsssb.delhi.gov.in/notifications",
				Selector:    ".notification, .notice, tr",
				TitleSelect: ".title, .notice-title, td:first-child",
				LinkSelect:  "a",
			},
			{
				Name:        "HSSC",
				URL:         "https://hssc.gov.in/advertisements",
				Selector:    ".advertisement, .notice, tr",
				TitleSelect: ".title, .adv-title, td:first-child",
				LinkSelect:  "a",
			},
			{
				Name:        "KPSC",
				URL:         "https://kpsc.kar.nic.in/notifications",
				Selector:    ".notification, .notice, tr",
				TitleSelect: ".title, .notice-title, td:first-child",
				LinkSelect:  "a",
			},
			{
				Name:        "TNPSC",
				URL:         "https://tnpsc.gov.in/notifications",
				Selector:    ".notification, .notice, tr",
				TitleSelect: ".title, .notice-title, td:first-child",
				LinkSelect:  "a",
			},
			{
				Name:        "WBPSC",
				URL:         "https://psc.wb.gov.in/notifications",
				Selector:    ".notification, .notice, tr",
				TitleSelect: ".title, .notice-title, td:first-child",
				LinkSelect:  "a",
			},
		},
	}
}

// Fetch retrieves jobs from all configured state commissions.
func (s *StatePSCSource) Fetch(ctx context.Context) ([]GovtJobSource, error) {
	log.Info().Msg("Starting crawl for source: State PSCs")

	var allJobs []GovtJobSource
	for _, state := range s.states {
		jobs, err := s.fetchState(ctx, state)
		if err != nil {
			log.Warn().Err(err).Str("state", state.Name).Msg("State PSC fetch failed")
			continue
		}
		allJobs = append(allJobs, jobs...)
	}

	log.Info().Int("totalJobs", len(allJobs)).Msg("State PSCs fetch completed")
	return allJobs, nil
}

func (s *StatePSCSource) fetchState(ctx context.Context, state stateConfig) ([]GovtJobSource, error) {
	// Rate limiting: 2 second delay between requests
	time.Sleep(2 * time.Second)

	req, err := http.NewRequestWithContext(ctx, "GET", state.URL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for %s: %w", state.Name, err)
	}
	req.Header.Set("User-Agent", "RojgarSetu/2.0")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch %s: %w", state.Name, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s returned status: %d", state.Name, resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse %s HTML: %w", state.Name, err)
	}

	var jobs []GovtJobSource
	doc.Find(state.Selector).Each(func(i int, item *goquery.Selection) {
		title := strings.TrimSpace(item.Find(state.TitleSelect).Text())
		link, exists := item.Find(state.LinkSelect).Attr("href")
		if !exists || title == "" {
			return
		}

		// Convert relative URLs to absolute
		if strings.HasPrefix(link, "/") {
			baseURL := strings.TrimSuffix(state.URL, "/")
			baseURL = baseURL[:strings.LastIndex(baseURL, "/")]
			link = baseURL + link
		}

		job := GovtJobSource{
			Source:           strings.ToLower(state.Name),
			Title:            title,
			Organization:     state.Name,
			ApplyURL:         link,
			LegalAttribution: fmt.Sprintf("Source: Official Portal (%s)", extractDomain(state.URL)),
			PostedDate:       time.Now(),
		}

		// SHA-256 hash for deduplication
		hash := sha256.Sum256([]byte(title + link))
		job.SHA256Hash = hex.EncodeToString(hash[:])

		if isValidGovtJob(&job) {
			jobs = append(jobs, job)
		}
	})

	log.Info().Int("jobs", len(jobs)).Str("state", state.Name).Msg("State PSC fetch successful")
	return jobs, nil
}

// extractDomain extracts the domain from a URL for legal attribution
func extractDomain(url string) string {
	if strings.HasPrefix(url, "http://") {
		url = strings.TrimPrefix(url, "http://")
	} else if strings.HasPrefix(url, "https://") {
		url = strings.TrimPrefix(url, "https://")
	}

	if idx := strings.Index(url, "/"); idx != -1 {
		url = url[:idx]
	}

	return url
}

// Name returns the source name.
func (s *StatePSCSource) Name() string {
	return s.NameStr
}

// FetchJobs implements the JobSource interface for compatibility with the engine
func (s *StatePSCSource) FetchJobs() ([]Job, error) {
	ctx := context.Background()
	govtJobs, err := s.Fetch(ctx)
	if err != nil {
		return nil, err
	}

	var jobs []Job
	for _, govtJob := range govtJobs {
		job := Job{
			Title:             govtJob.Title,
			CompanyOrDept:     govtJob.Organization,
			Location:          "",
			QualificationReq:  "",
			SalaryOrPayScale:  "",
			ApplyURL:          govtJob.ApplyURL,
			SourceAttribution: govtJob.LegalAttribution,
			HashChecksum:      govtJob.SHA256Hash,
		}
		jobs = append(jobs, job)
	}

	return jobs, nil
}
