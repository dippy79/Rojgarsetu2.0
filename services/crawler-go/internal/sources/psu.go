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

// PSUSource fetches jobs from Public Sector Undertakings and Scientific Research Organizations.
// It supports: ISRO, DRDO, BARC, ONGC, NTPC, BHEL, and IOCL.
type PSUSource struct {
	BaseSource
	httpClient    *http.Client
	organizations []psuOrganization
}

// psuOrganization holds configuration for each PSU organization
type psuOrganization struct {
	Name        string
	URL         string
	Selector    string // CSS selector for job listings
	TitleSelect string // CSS selector for job title
	LinkSelect  string // CSS selector for job link
}

// NewPSUSource creates a PSU source with multi-organization support.
func NewPSUSource() *PSUSource {
	return &PSUSource{
		BaseSource: BaseSource{
			NameStr: "psu",
			BaseURL: "",
		},
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		organizations: []psuOrganization{
			{
				Name:        "ISRO",
				URL:         "https://www.isro.gov.in/Careers.html",
				Selector:    ".career-item, .notice-item, .advertisement, tr",
				TitleSelect: ".title, .notice-title, .adv-title, td:first-child",
				LinkSelect:  "a",
			},
			{
				Name:        "DRDO",
				URL:         "https://drdo.gov.in/careers",
				Selector:    ".career-item, .notice-item, .advertisement, tr",
				TitleSelect: ".title, .notice-title, .adv-title, td:first-child",
				LinkSelect:  "a",
			},
			{
				Name:        "BARC",
				URL:         "https://barc.gov.in/careers",
				Selector:    ".career-item, .notice-item, .advertisement, tr",
				TitleSelect: ".title, .notice-title, .adv-title, td:first-child",
				LinkSelect:  "a",
			},
			{
				Name:        "ONGC",
				URL:         "https://ongcindia.com/careers",
				Selector:    ".career-item, .notice-item, .advertisement, tr",
				TitleSelect: ".title, .notice-title, .adv-title, td:first-child",
				LinkSelect:  "a",
			},
			{
				Name:        "NTPC",
				URL:         "https://careers.ntpc.co.in",
				Selector:    ".career-item, .notice-item, .advertisement, tr",
				TitleSelect: ".title, .notice-title, .adv-title, td:first-child",
				LinkSelect:  "a",
			},
			{
				Name:        "BHEL",
				URL:         "https://bhel.com/careers",
				Selector:    ".career-item, .notice-item, .advertisement, tr",
				TitleSelect: ".title, .notice-title, .adv-title, td:first-child",
				LinkSelect:  "a",
			},
			{
				Name:        "IOCL",
				URL:         "https://iocl.com/apprenticeships",
				Selector:    ".career-item, .notice-item, .advertisement, tr",
				TitleSelect: ".title, .notice-title, .adv-title, td:first-child",
				LinkSelect:  "a",
			},
		},
	}
}

// Fetch retrieves jobs from all configured PSU organizations.
func (s *PSUSource) Fetch(ctx context.Context) ([]GovtJobSource, error) {
	log.Info().Msg("Starting crawl for source: PSU")

	var allJobs []GovtJobSource
	for _, org := range s.organizations {
		jobs, err := s.fetchOrganization(ctx, org)
		if err != nil {
			log.Warn().Err(err).Str("organization", org.Name).Msg("PSU organization fetch failed")
			continue
		}
		allJobs = append(allJobs, jobs...)
	}

	log.Info().Int("totalJobs", len(allJobs)).Msg("PSU fetch completed")
	return allJobs, nil
}

func (s *PSUSource) fetchOrganization(ctx context.Context, org psuOrganization) ([]GovtJobSource, error) {
	// Rate limiting: 2 second delay between requests
	time.Sleep(2 * time.Second)

	req, err := http.NewRequestWithContext(ctx, "GET", org.URL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for %s: %w", org.Name, err)
	}
	req.Header.Set("User-Agent", "RojgarSetu/2.0")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch %s: %w", org.Name, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s returned status: %d", org.Name, resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse %s HTML: %w", org.Name, err)
	}

	var jobs []GovtJobSource
	doc.Find(org.Selector).Each(func(i int, item *goquery.Selection) {
		title := strings.TrimSpace(item.Find(org.TitleSelect).Text())
		link, exists := item.Find(org.LinkSelect).Attr("href")
		if !exists || title == "" {
			return
		}

		// Convert relative URLs to absolute
		if strings.HasPrefix(link, "/") {
			baseURL := strings.TrimSuffix(org.URL, "/")
			link = baseURL + link
		}

		job := GovtJobSource{
			Source:           strings.ToLower(org.Name),
			Title:            title,
			Organization:     org.Name,
			ApplyURL:         link,
			LegalAttribution: fmt.Sprintf("Source: Official Portal (%s)", extractDomain(org.URL)),
			PostedDate:       time.Now(),
		}

		// SHA-256 hash for deduplication
		hash := sha256.Sum256([]byte(title + link))
		job.SHA256Hash = hex.EncodeToString(hash[:])

		if isValidGovtJob(&job) {
			jobs = append(jobs, job)
		}
	})

	log.Info().Int("jobs", len(jobs)).Str("organization", org.Name).Msg("PSU organization fetch successful")
	return jobs, nil
}

// Name returns the source name.
func (s *PSUSource) Name() string {
	return s.NameStr
}

// FetchJobs implements the JobSource interface for compatibility with the engine
func (s *PSUSource) FetchJobs() ([]Job, error) {
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
