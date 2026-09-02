package courses

import (
	"github.com/rojgarsetu/crawler/internal/shared"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

// NSDCSource scrapes courses from NSDC (National Skill Development Corporation)
type NSDCSource struct {
	shared.BaseSource
	client *http.Client
}

// NewNSDCSource creates a new NSDC source
func NewNSDCSource() *NSDCSource {
	return &NSDCSource{
		BaseSource: shared.BaseSource{NameStr: "nsdc", BaseURL: "https://www.nsdcindia.org"},
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Fetch retrieves courses from NSDC
func (s *NSDCSource) Fetch(ctx context.Context) ([]shared.CourseSource, error) {
	log.Info().Msg("Starting crawl for source: NSDC")

	var courses []shared.CourseSource

	nsdcCourses, err := s.fetchFromWebsite(ctx)
	if err != nil {
		log.Error().Err(err).Msg("NSDC fetch failed")
		return nil, fmt.Errorf("failed to fetch from NSDC: %w", err)
	}

	courses = append(courses, nsdcCourses...)
	log.Info().Int("totalCourses", len(courses)).Msg("NSDC fetch completed")
	return courses, nil
}

// fetchFromWebsite fetches from NSDC website
func (s *NSDCSource) fetchFromWebsite(ctx context.Context) ([]shared.CourseSource, error) {
	urls := []string{
		"https://www.nsdcindia.org/skill-training",
		"https://www.nsdcindia.org/courses",
		"https://www.skillindia.gov.in/courses",
	}

	var allCourses []shared.CourseSource

	for _, url := range urls {
		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
		if err != nil {
			continue
		}

		req.Header.Set("User-Agent", "RojgarSetu/2.0")

		resp, err := s.client.Do(req)
		if err != nil {
			continue
		}

		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			continue
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			continue
		}

		courses := s.parseHTMLCourses(string(body), url)
		allCourses = append(allCourses, courses...)
	}

	log.Info().Int("coursesFromWebsite", len(allCourses)).Msg("NSDC website fetch successful")
	return allCourses, nil
}

// parseHTMLCourses parses courses from NSDC HTML
func (s *NSDCSource) parseHTMLCourses(html, baseURL string) []shared.CourseSource {
	var courses []shared.CourseSource

	patterns := []string{
		`<a[^>]*href="(/[^"]*course[^"]*)"[^>]*>.*?<h[34][^>]*>([^<]*)</h[34]>`,
		`<div[^>]*class="[^"]*course[^"]*"[^>]*>.*?<h3[^>]*>([^<]*)</h3>.*?<a[^>]*href="([^"]*)"`,
		`<a[^>]*href="(/[^"]*training[^"]*)"[^>]*>([^<]*)</a>`,
	}

	for _, pattern := range patterns {
		matches := shared.ExtractMatches(html, pattern)
		for _, match := range matches {
			if len(match) >= 3 {
				link := strings.TrimSpace(match[1])
				title := strings.TrimSpace(match[2])

				if len(title) > 5 {
					// Make absolute URL
					if !strings.HasPrefix(link, "http") {
						link = "https://www.nsdcindia.org" + link
					}

					course := shared.CourseSource{
						Source:    "nsdc",
						Provider:  "NSDC",
						Title:     title,
						URL:       link,
						IsFree:    false,
						Mode:      "offline",
						CreatedAt: time.Now(),
					}
					if shared.IsValidCourse(&course) {
						courses = append(courses, course)
					}
				}
			}
		}
	}

	return courses
}

// Name returns the source name
func (s *NSDCSource) Name() string {
	return s.NameStr
}


