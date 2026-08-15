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

// DefenceSource fetches jobs from Defense and Security Forces.
// It supports: Indian Army, Air Force (AFCAT), Navy, and Coast Guard.
type DefenceSource struct {
	BaseSource
	httpClient *http.Client
	services   []defenceService
}

// defenceService holds configuration for each defense service
type defenceService struct {
	Name        string
	URL         string
	Selector    string // CSS selector for job listings
	TitleSelect string // CSS selector for job title
	LinkSelect  string // CSS selector for job link
}

// NewDefenceSource creates a Defense source with multi-service support.
func NewDefenceSource() *DefenceSource {
	return &DefenceSource{
		BaseSource: BaseSource{
			NameStr: "defence",
			BaseURL: "",
		},
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		services: []defenceService{
			{
				Name:        "Indian Army",
				URL:         "https://joinindianarmy.nic.in",
				Selector:    ".notice-item, .notification, .advertisement, tr",
				TitleSelect: ".title, .notice-title, .adv-title, td:first-child",
				LinkSelect:  "a",
			},
			{
				Name:        "Indian Air Force (AFCAT)",
				URL:         "https://afcat.cdac.in",
				Selector:    ".notice-item, .notification, .advertisement, tr",
				TitleSelect: ".title, .notice-title, .adv-title, td:first-child",
				LinkSelect:  "a",
			},
			{
				Name:        "Indian Navy",
				URL:         "https://joinindiannavy.gov.in",
				Selector:    ".notice-item, .notification, .advertisement, tr",
				TitleSelect: ".title, .notice-title, .adv-title, td:first-child",
				LinkSelect:  "a",
			},
			{
				Name:        "Indian Coast Guard",
				URL:         "https://joinindiancoastguard.cdac.in",
				Selector:    ".notice-item, .notification, .advertisement, tr",
				TitleSelect: ".title, .notice-title, .adv-title, td:first-child",
				LinkSelect:  "a",
			},
		},
	}
}

// Fetch retrieves jobs from all configured defense services.
func (s *DefenceSource) Fetch(ctx context.Context) ([]GovtJobSource, error) {
	log.Info().Msg("Starting crawl for source: Defense")

	var allJobs []GovtJobSource
	for _, service := range s.services {
		jobs, err := s.fetchService(ctx, service)
		if err != nil {
			log.Warn().Err(err).Str("service", service.Name).Msg("Defense service fetch failed")
			continue
		}
		allJobs = append(allJobs, jobs...)
	}

	log.Info().Int("totalJobs", len(allJobs)).Msg("Defense fetch completed")
	return allJobs, nil
}

func (s *DefenceSource) fetchService(ctx context.Context, service defenceService) ([]GovtJobSource, error) {
	// Rate limiting: 2 second delay between requests
	time.Sleep(2 * time.Second)

	req, err := http.NewRequestWithContext(ctx, "GET", service.URL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for %s: %w", service.Name, err)
	}
	req.Header.Set("User-Agent", "RojgarSetu/2.0")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch %s: %w", service.Name, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s returned status: %d", service.Name, resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse %s HTML: %w", service.Name, err)
	}

	var jobs []GovtJobSource
	doc.Find(service.Selector).Each(func(i int, item *goquery.Selection) {
		title := strings.TrimSpace(item.Find(service.TitleSelect).Text())
		link, exists := item.Find(service.LinkSelect).Attr("href")
		if !exists || title == "" {
			return
		}

		// Convert relative URLs to absolute
		if strings.HasPrefix(link, "/") {
			baseURL := strings.TrimSuffix(service.URL, "/")
			link = baseURL + link
		}

		job := GovtJobSource{
			Source:           strings.ToLower(strings.ReplaceAll(service.Name, " ", "_")),
			Title:            title,
			Organization:     service.Name,
			ApplyURL:         link,
			LegalAttribution: fmt.Sprintf("Source: Official Portal (%s)", extractDomain(service.URL)),
			PostedDate:       time.Now(),
		}

		// SHA-256 hash for deduplication
		hash := sha256.Sum256([]byte(title + link))
		job.SHA256Hash = hex.EncodeToString(hash[:])

		if isValidGovtJob(&job) {
			jobs = append(jobs, job)
		}
	})

	log.Info().Int("jobs", len(jobs)).Str("service", service.Name).Msg("Defense service fetch successful")
	return jobs, nil
}

// Name returns the source name.
func (s *DefenceSource) Name() string {
	return s.NameStr
}

// FetchJobs implements the JobSource interface for compatibility with the engine
func (s *DefenceSource) FetchJobs() ([]Job, error) {
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
