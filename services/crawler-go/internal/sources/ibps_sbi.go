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

// IbpsSbiSource fetches banking job notices from IBPS and SBI Careers.
// It uses polite HTML scraping with legal attribution and SHA-256 deduplication.
type IbpsSbiSource struct {
	BaseSource
	httpClient *http.Client
}

// NewIbpsSbiSource creates an IBPS/SBI source.
func NewIbpsSbiSource() *IbpsSbiSource {
	return &IbpsSbiSource{
		BaseSource: BaseSource{
			NameStr: "ibps_sbi",
			BaseURL: "https://www.ibps.in",
		},
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Fetch retrieves banking jobs from IBPS and SBI.
func (s *IbpsSbiSource) Fetch(ctx context.Context) ([]GovtJobSource, error) {
	log.Info().Msg("Starting crawl for source: IBPS/SBI")

	var allJobs []GovtJobSource

	// Fetch from IBPS
	ibpsJobs, err := s.fetchIBPS(ctx)
	if err != nil {
		log.Warn().Err(err).Msg("IBPS fetch failed, continuing with SBI")
	} else {
		allJobs = append(allJobs, ibpsJobs...)
	}

	// Fetch from SBI
	sbiJobs, err := s.fetchSBI(ctx)
	if err != nil {
		log.Warn().Err(err).Msg("SBI fetch failed")
	} else {
		allJobs = append(allJobs, sbiJobs...)
	}

	log.Info().Int("totalJobs", len(allJobs)).Msg("IBPS/SBI fetch completed")
	return allJobs, nil
}

func (s *IbpsSbiSource) fetchIBPS(ctx context.Context) ([]GovtJobSource, error) {
	url := "https://www.ibps.in/career.htm"

	// Rate limiting: 2 second delay
	time.Sleep(2 * time.Second)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for IBPS: %w", err)
	}
	req.Header.Set("User-Agent", "RojgarSetu/2.0")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch IBPS: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("IBPS returned status: %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse IBPS HTML: %w", err)
	}

	var jobs []GovtJobSource
	doc.Find("table tbody tr").Each(func(i int, row *goquery.Selection) {
		cols := row.Find("td")
		if cols.Length() < 2 {
			return
		}

		title := strings.TrimSpace(cols.Eq(0).Text())
		link, exists := cols.Eq(1).Find("a").Attr("href")
		if !exists || title == "" {
			return
		}

		// Convert relative URLs to absolute
		if strings.HasPrefix(link, "/") {
			link = "https://www.ibps.in" + link
		}

		job := GovtJobSource{
			Source:           "ibps",
			Title:            title,
			Organization:     "IBPS",
			ApplyURL:         link,
			LegalAttribution: "Source: Official Portal (ibps.in)",
			PostedDate:       time.Now(),
		}

		// SHA-256 hash for deduplication
		hash := sha256.Sum256([]byte(title + link))
		job.SHA256Hash = hex.EncodeToString(hash[:])

		if isValidGovtJob(&job) {
			jobs = append(jobs, job)
		}
	})

	log.Info().Int("jobs", len(jobs)).Msg("IBPS fetch successful")
	return jobs, nil
}

func (s *IbpsSbiSource) fetchSBI(ctx context.Context) ([]GovtJobSource, error) {
	url := "https://sbi.co.in/web/careers"

	// Rate limiting: 2 second delay
	time.Sleep(2 * time.Second)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for SBI: %w", err)
	}
	req.Header.Set("User-Agent", "RojgarSetu/2.0")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch SBI: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("SBI returned status: %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse SBI HTML: %w", err)
	}

	var jobs []GovtJobSource
	doc.Find(".career-item, .notice-item, .latest-notice").Each(func(i int, item *goquery.Selection) {
		title := strings.TrimSpace(item.Find("h3, .title, .notice-title").Text())
		link, exists := item.Find("a").Attr("href")
		if !exists || title == "" {
			return
		}

		// Convert relative URLs to absolute
		if strings.HasPrefix(link, "/") {
			link = "https://sbi.co.in" + link
		}

		job := GovtJobSource{
			Source:           "sbi",
			Title:            title,
			Organization:     "State Bank of India",
			ApplyURL:         link,
			LegalAttribution: "Source: Official Portal (sbi.co.in)",
			PostedDate:       time.Now(),
		}

		// SHA-256 hash for deduplication
		hash := sha256.Sum256([]byte(title + link))
		job.SHA256Hash = hex.EncodeToString(hash[:])

		if isValidGovtJob(&job) {
			jobs = append(jobs, job)
		}
	})

	log.Info().Int("jobs", len(jobs)).Msg("SBI fetch successful")
	return jobs, nil
}

// Name returns the source name.
func (s *IbpsSbiSource) Name() string {
	return s.NameStr
}

// FetchJobs implements the JobSource interface for compatibility with the engine
func (s *IbpsSbiSource) FetchJobs() ([]Job, error) {
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
