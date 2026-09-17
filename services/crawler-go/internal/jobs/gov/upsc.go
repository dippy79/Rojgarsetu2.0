package gov

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/rojgarsetu/crawler/internal/browser"
	"github.com/rojgarsetu/crawler/internal/sources"
)

type UPSCSource struct {
	sources.BaseSource
	client *http.Client
}

func NewUPSCSource(pool *browser.Pool) *UPSCSource {
	return &UPSCSource{
		BaseSource: sources.BaseSource{NameStr: "upsc", BaseURL: "https://www.upsc.gov.in"},
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (s *UPSCSource) Fetch(ctx context.Context) ([]sources.GovJobSource, error) {
	log.Println("[UPSC] Fetching live jobs from upsc.gov.in...")

	url := "https://www.upsc.gov.in/examinations/active-exams"
	res, err := s.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	var jobs []sources.GovJobSource
	doc.Find("table tbody tr").Each(func(i int, row *goquery.Selection) {
		title := strings.TrimSpace(row.Find("td").Eq(1).Text())
		if title == "" {
			return
		}

		link, _ := row.Find("td").Eq(1).Find("a").Attr("href")
		if link != "" && !strings.HasPrefix(link, "http") {
			link = "https://www.upsc.gov.in" + link
		}

		job := sources.GovJobSource{
			Source:       "upsc",
			Title:        title,
			Department:   "Union Public Service Commission",
			ApplyURL:     link,
			Category:     "CENTRAL",
			StateName:    "ALL_INDIA",
			CreatedAt:    time.Now(),
		}
		jobs = append(jobs, job)
	})

	log.Printf("[UPSC] Successfully extracted %d live listings", len(jobs))
	return jobs, nil
}

func (s *UPSCSource) Name() string {
	return s.NameStr
}
