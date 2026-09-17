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

type SSCSource struct {
	sources.BaseSource
	client *http.Client
}

func NewSSCSource(pool *browser.Pool) *SSCSource {
	return &SSCSource{
		BaseSource: sources.BaseSource{NameStr: "ssc", BaseURL: "https://ssc.gov.in"},
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (s *SSCSource) Fetch(ctx context.Context) ([]sources.GovJobSource, error) {
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

	var jobs []sources.GovJobSource
	doc.Find("table tr, .notice-item").Each(func(i int, item *goquery.Selection) {
		title := strings.TrimSpace(item.Text())
		if len(title) < 10 {
			return
		}

		link, _ := item.Find("a").Attr("href")
		if link != "" && !strings.HasPrefix(link, "http") {
			link = "https://ssc.gov.in" + link
		}

		job := sources.GovJobSource{
			Source:       "ssc",
			Title:        title,
			Department:   "Staff Selection Commission",
			ApplyURL:     link,
			Category:     "CENTRAL",
			StateName:    "ALL_INDIA",
			CreatedAt:    time.Now(),
		}
		jobs = append(jobs, job)
	})

	log.Printf("[SSC] Successfully extracted %d live listings", len(jobs))
	return jobs, nil
}

func (s *SSCSource) Name() string {
	return s.NameStr
}
