package sources

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

// CourseraSource scrapes courses from Coursera
type CourseraSource struct {
	BaseSource
	client *http.Client
	apiURL string
}

// NewCourseraSource creates a new Coursera source
func NewCourseraSource() *CourseraSource {
	return &CourseraSource{
		BaseSource: BaseSource{NameStr: "coursera", BaseURL: "https://www.coursera.org"},
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		apiURL: "https://api.coursera.org/api/courses.v1",
	}
}

// Fetch retrieves courses from Coursera
func (s *CourseraSource) Fetch(ctx context.Context) ([]CourseSource, error) {
	log.Info().Msg("Starting crawl for source: Coursera")

	var courses []CourseSource

	courseraCourses, err := s.fetchFromAPI(ctx)
	if err != nil {
		log.Warn().Err(err).Msg("Coursera API fetch failed, trying website")
		courseraCourses, err = s.fetchFromWebsite(ctx)
		if err != nil {
			log.Error().Err(err).Msg("All Coursera fetch methods failed")
			return nil, fmt.Errorf("failed to fetch from Coursera: %w", err)
		}
	}

	courses = append(courses, courseraCourses...)
	log.Info().Int("totalCourses", len(courses)).Msg("Coursera fetch completed")
	return courses, nil
}

// fetchFromAPI fetches from Coursera API
func (s *CourseraSource) fetchFromAPI(ctx context.Context) ([]CourseSource, error) {
	// Coursera has a public API but requires authentication for full access
	// Try with limited parameters
	req, err := http.NewRequestWithContext(ctx, "GET", s.apiURL+"?fields=description,photoUrl,slug,startDate,level", nil)
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
		return nil, fmt.Errorf("Coursera API returned status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var courseraData struct {
		Elements []struct {
			ID          string `json:"id"`
			Name        string `json:"name"`
			Description string `json:"description"`
			Slug        string `json:"slug"`
			Level       string `json:"level"`
			PhotoURL    string `json:"photoUrl"`
			StartDate   string `json:"startDate"`
		} `json:"elements"`
	}

	if err := json.Unmarshal(body, &courseraData); err != nil {
		return nil, err
	}

	var courses []CourseSource
	for _, c := range courseraData.Elements {
		course := CourseSource{
			Source:       "coursera",
			Provider:     "Coursera",
			Title:        cleanString(c.Name),
			URL:          "https://www.coursera.org/learn/" + c.Slug,
			Level:        normalizeCourseLevel(c.Level),
			Description:  cleanString(c.Description),
			ThumbnailURL: c.PhotoURL,
			IsFree:       false,
			CreatedAt:    time.Now(),
		}

		if course.Title != "" && isValidCourse(&course) {
			courses = append(courses, course)
		}
	}

	log.Info().Int("coursesFromAPI", len(courses)).Msg("Coursera API fetch successful")
	return courses, nil
}

// fetchFromWebsite fetches from Coursera website
func (s *CourseraSource) fetchFromWebsite(ctx context.Context) ([]CourseSource, error) {
	urls := []string{
		"https://www.coursera.org/browse",
		"https://www.coursera.org/courses",
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

	log.Info().Int("coursesFromWebsite", len(allCourses)).Msg("Coursera website fetch successful")
	return allCourses, nil
}

// parseHTMLCourses parses courses from Coursera HTML
func (s *CourseraSource) parseHTMLCourses(html string) []CourseSource {
	var courses []CourseSource

	// Look for JSON-LD structured data
	pattern := `<script[^>]*type="application/ld\+json"[^>]*>([^<]+)</script>`
	matches := extractMatches(html, pattern)

	for _, match := range matches {
		if len(match) >= 2 {
			jsonData := match[1]

			var courseData struct {
				Name        string `json:"name"`
				Description string `json:"description"`
				URL         string `json:"url"`
				Level       string `json:"level"`
			}

			if err := json.Unmarshal([]byte(jsonData), &courseData); err != nil {
				continue
			}

			if courseData.Name != "" && courseData.URL != "" {
				course := CourseSource{
					Source:      "coursera",
					Provider:    "Coursera",
					Title:       cleanString(courseData.Name),
					URL:         courseData.URL,
					Level:       normalizeCourseLevel(courseData.Level),
					Description: cleanString(courseData.Description),
					IsFree:      false,
					CreatedAt:   time.Now(),
				}
				if isValidCourse(&course) {
					courses = append(courses, course)
				}
			}
		}
	}

	return courses
}

// Name returns the source name
func (s *CourseraSource) Name() string {
	return s.NameStr
}
