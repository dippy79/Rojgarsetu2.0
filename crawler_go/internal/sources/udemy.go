package sources

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

// UdemySource scrapes courses from Udemy
type UdemySource struct {
	BaseSource
	client *http.Client
}

// NewUdemySource creates a new Udemy source
func NewUdemySource() *UdemySource {
	return &UdemySource{
		BaseSource: BaseSource{NameStr: "udemy", BaseURL: "https://www.udemy.com"},
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Fetch retrieves courses from Udemy
func (s *UdemySource) Fetch(ctx context.Context) ([]CourseSource, error) {
	log.Info().Msg("Starting crawl for source: Udemy")

	var courses []CourseSource

	udemyCourses, err := s.fetchFromWebsite(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Udemy fetch failed")
		return nil, fmt.Errorf("failed to fetch from Udemy: %w", err)
	}

	courses = append(courses, udemyCourses...)
	log.Info().Int("totalCourses", len(courses)).Msg("Udemy fetch completed")
	return courses, nil
}

// fetchFromWebsite fetches from Udemy website
func (s *UdemySource) fetchFromWebsite(ctx context.Context) ([]CourseSource, error) {
	urls := []string{
		"https://www.udemy.com/tech/",
		"https://www.udemy.com/business/",
		"https://www.udemy.com/academics/",
	}

	var allCourses []CourseSource

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

		courses := s.parseHTMLCourses(string(body))
		allCourses = append(allCourses, courses...)
	}

	log.Info().Int("coursesFromWebsite", len(allCourses)).Msg("Udemy website fetch successful")
	return allCourses, nil
}

// parseHTMLCourses parses courses from Udemy HTML
func (s *UdemySource) parseHTMLCourses(html string) []CourseSource {
	var courses []CourseSource

	// Look for course cards
	patterns := []string{
		`"title":"([^"]+)".*?"url":"([^"]+)".*?"price":"([^"]+)".*?"duration":"([^"]+)"`,
		`<a[^>]*href="(/course/[^"]*)"[^>]*>.*?<h3[^>]*>([^<]*)</h3>`,
	}

	for _, pattern := range patterns {
		matches := extractMatches(html, pattern)
		for _, match := range matches {
			if len(match) >= 3 {
				url := strings.TrimSpace(match[1])
				title := strings.TrimSpace(match[2])

				if len(title) > 5 {
					// Make absolute URL
					if !strings.HasPrefix(url, "http") {
						url = "https://www.udemy.com" + url
					}

					course := CourseSource{
						Source:    "udemy",
						Provider:  "Udemy",
						Title:     title,
						URL:       url,
						IsFree:    false,
						CreatedAt: time.Now(),
					}
					if isValidCourse(&course) {
						courses = append(courses, course)
					}
				}
			}
		}
	}

	return courses
}

// Name returns the source name
func (s *UdemySource) Name() string {
	return s.NameStr
}
