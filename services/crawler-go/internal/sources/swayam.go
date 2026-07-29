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

// SWAYAMSource scrapes courses from SWAYAM
type SWAYAMSource struct {
	BaseSource
	client *http.Client
	apiURL string
}

// NewSWAYAMSource creates a new SWAYAM source
func NewSWAYAMSource() *SWAYAMSource {
	return &SWAYAMSource{
		BaseSource: BaseSource{NameStr: "swayam", BaseURL: "https://swayam.gov.in"},
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		apiURL: "https://swayam.gov.in/api/courses",
	}
}

// Fetch retrieves courses from SWAYAM
func (s *SWAYAMSource) Fetch(ctx context.Context) ([]CourseSource, error) {
	log.Info().Msg("Starting crawl for source: SWAYAM")

	var courses []CourseSource

	apiCourses, err := s.fetchFromAPI(ctx)
	if err != nil {
		log.Warn().Err(err).Msg("SWAYAM API fetch failed, trying website")
		apiCourses, err = s.fetchFromWebsite(ctx)
		if err != nil {
			log.Error().Err(err).Msg("All SWAYAM fetch methods failed")
			return nil, fmt.Errorf("failed to fetch from SWAYAM: %w", err)
		}
	}

	courses = append(courses, apiCourses...)
	log.Info().Int("totalCourses", len(courses)).Msg("SWAYAM fetch completed")
	return courses, nil
}

// fetchFromAPI fetches from SWAYAM API
func (s *SWAYAMSource) fetchFromAPI(ctx context.Context) ([]CourseSource, error) {
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
		return nil, fmt.Errorf("SWAYAM API returned status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var swayamData struct {
		Courses []struct {
			Title       string `json:"title"`
			Description string `json:"description"`
			Duration    string `json:"duration"`
			Level       string `json:"level"`
			URL         string `json:"url"`
			Thumbnail   string `json:"thumbnail"`
			StartDate   string `json:"startDate"`
		} `json:"courses"`
	}

	if err := json.Unmarshal(body, &swayamData); err != nil {
		return nil, err
	}

	var courses []CourseSource
	for _, c := range swayamData.Courses {
		course := CourseSource{
			Source:       "swayam",
			Provider:     "SWAYAM",
			Title:        cleanString(c.Title),
			URL:          "https://swayam.gov.in" + c.URL,
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

	log.Info().Int("coursesFromAPI", len(courses)).Msg("SWAYAM API fetch successful")
	return courses, nil
}

// fetchFromWebsite fetches from SWAYAM website
func (s *SWAYAMSource) fetchFromWebsite(ctx context.Context) ([]CourseSource, error) {
	urls := []string{
		"https://swayam.gov.in/explore",
		"https://swayam.gov.in/courses",
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

	log.Info().Int("coursesFromWebsite", len(allCourses)).Msg("SWAYAM website fetch successful")
	return allCourses, nil
}

// parseHTMLCourses parses courses from SWAYAM HTML
func (s *SWAYAMSource) parseHTMLCourses(html string) []CourseSource {
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
						Source:    "swayam",
						Provider:  "SWAYAM",
						Title:     title,
						URL:       "https://swayam.gov.in" + link,
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
func (s *SWAYAMSource) Name() string {
	return s.NameStr
}
