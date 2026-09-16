package sources

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// SSCScraper scrapes SSC official portal
type SSCScraper struct {
	client *http.Client
}

// NewSSCScraper creates a new SSC scraper
func NewSSCScraper(client any) *SSCScraper {
	return &SSCScraper{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// FetchJobs fetches real jobs from SSC Notices page
func (s *SSCScraper) FetchJobs() ([]Job, error) {
	log.Println("[SSC] Fetching live jobs from ssc.gov.in...")

	url := "https://ssc.gov.in/Notices"
	res, err := s.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	var jobs []Job
	// SSC notices usually in a list or table
	doc.Find(".notice-item, tr").Each(func(i int, item *goquery.Selection) {
		title := strings.TrimSpace(item.Text())
		if len(title) < 10 {
			return
		}

		link, _ := item.Find("a").Attr("href")
		if link == "" {
			return
		}

		job := Job{
			Title:             title,
			CompanyOrDept:     "Staff Selection Commission",
			Location:          "All India",
			QualificationReq:  "10th / 12th / Graduate",
			SalaryOrPayScale:  "As per 7th CPC",
			ApplyURL:          link,
			SourceAttribution: "Source: SSC Official Portal (ssc.gov.in)",
		}
		jobs = append(jobs, job)
	})

	if len(jobs) == 0 {
		log.Println("[SSC] Warning: No live jobs found via parser.")
	}

	log.Printf("[SSC] Successfully extracted %d live listings", len(jobs))
	return jobs, nil
}
