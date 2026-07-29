package sources

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

// NPTELSource scrapes courses from NPTEL
type NPTELSource struct {
	BaseSource
	client *http.Client
	apiURL string
}

// NewNPTELSource creates a new NPTEL source
func NewNPTELSource() *NPTELSource {
	return &NPTELSource{
		BaseSource: BaseSource{NameStr: "nptel", BaseURL: "https://nptel.ac.in"},
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		apiURL: "https://nptel.ac.in/api/course.list",
	}
}

// Fetch retrieves courses from NPTEL
func (s *NPTELSource) Fetch(ctx context.Context) ([]CourseSource, error) {
	log.Info().Msg("Starting crawl for source: NPTEL")

	var courses []CourseSource

	// Try API first
	apiCourses, err := s.fetchFromAPI(ctx)
	if err != nil {
		log.Warn().Err(err).Msg("NPTEL API fetch failed, trying website")
		apiCourses, err = s.fetchFromWebsite(ctx)
		if err != nil {
			log.Error().Err(err).Msg("All NPTEL fetch methods failed")
			return nil, fmt.Errorf("failed to fetch from NPTEL: %w", err)
		}
	}

	courses = append(courses, apiCourses...)
	log.Info().Int("totalCourses", len(courses)).Msg("NPTEL fetch completed")
	return courses, nil
}

// fetchFromAPI fetches from NPTEL API
func (s *NPTELSource) fetchFromAPI(ctx context.Context) ([]CourseSource, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", s.apiURL, nil)
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
		return nil, fmt.Errorf("NPTEL API returned status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Try to parse JSON response
	var nptelData struct {
		Courses []struct {
			Title       string `json:"title"`
			Description string `json:"description"`
			Duration    string `json:"duration"`
			Level       string `json:"level"`
			URL         string `json:"url"`
			Thumbnail   string `json:"thumbnail"`
		} `json:"courses"`
	}

	if err := json.Unmarshal(body, &nptelData); err != nil {
		return nil, err
	}

	var courses []CourseSource
	for _, c := range nptelData.Courses {
		course := CourseSource{
			Source:       "nptel",
			Provider:     "NPTEL",
			Title:        cleanString(c.Title),
			URL:          "https://nptel.ac.in" + c.URL,
			Duration:     c.Duration,
			Level:        normalizeCourseLevel(c.Level),
			Description:  cleanString(c.Description),
			ThumbnailURL: c.Thumbnail,
			IsFree:       true,
			CreatedAt:    time.Now(),
		}

		if course.Title != "" && isValidCourse(&course) {
			courses = append(courses, course)
		}
	}

	log.Info().Int("coursesFromAPI", len(courses)).Msg("NPTEL API fetch successful")
	return courses, nil
}

// fetchFromWebsite fetches from NPTEL website
func (s *NPTELSource) fetchFromWebsite(ctx context.Context) ([]CourseSource, error) {
	urls := []string{
		"https://nptel.ac.in/course",
		"https://nptel.ac.in/online-course",
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

	log.Info().Int("coursesFromWebsite", len(allCourses)).Msg("NPTEL website fetch successful")
	return allCourses, nil
}

// parseHTMLCourses parses courses from NPTEL HTML
func (s *NPTELSource) parseHTMLCourses(html string) []CourseSource {
	var courses []CourseSource

	patterns := []string{
		`<a[^>]*href="(/course/[^"]*)"[^>]*>.*?<h3[^>]*>([^<]*)</h3>`,
		`<div[^>]*class="[^"]*course[^"]*"[^>]*>.*?<a[^>]*href="([^"]*)"[^>]*>([^<]*)</a>`,
	}

	for _, pattern := range patterns {
		matches := extractMatches(html, pattern)
		for _, match := range matches {
			if len(match) >= 3 {
				link := strings.TrimSpace(match[1])
				title := strings.TrimSpace(match[2])

				if len(title) > 5 {
					course := CourseSource{
						Source:    "nptel",
						Provider:  "NPTEL",
						Title:     title,
						URL:       "https://nptel.ac.in" + link,
						IsFree:    true,
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
func (s *NPTELSource) Name() string {
	return s.NameStr
}
