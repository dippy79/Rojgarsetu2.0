package courses

import (
	"github.com/rojgarsetu/crawler/internal/shared"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

// CourseraSource scrapes courses from Coursera
type CourseraSource struct {
	shared.BaseSource
	client *http.Client
	apiURL string
}

// NewCourseraSource creates a new Coursera source
func NewCourseraSource() *CourseraSource {
	return &CourseraSource{
		BaseSource: shared.BaseSource{NameStr: "coursera", BaseURL: "https://www.coursera.org"},
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		apiURL: "https://api.coursera.org/api/courses.v1",
	}
}

// Fetch retrieves courses from Coursera
func (s *CourseraSource) Fetch(ctx context.Context) ([]shared.CourseSource, error) {
	log.Info().Msg("Starting crawl for source: Coursera")

	var courses []shared.CourseSource

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

// fetchFromAPI fetches from Coursera API.
//
// Verified live (response shape as of this fix):
//
//	GET https://api.coursera.org/api/courses.v1?limit=3&fields=name,slug,startDate,photoUrl,level
//	200 OK
//	{
//	  "elements": [
//	    {
//	      "courseType": "v2.ondemand",
//	      "photoUrl": "https://...",
//	      "id": "l31la3mKEe-zFg7heHyXOQ",
//	      "slug": "googlecloud-...",
//	      "level": "BEGINNER",
//	      "partnerIds": ["443"],
//	      "name": "...",
//	      "startDate": 1727880158838   // <-- Unix epoch MILLISECONDS (number), not a string
//	    }
//	  ],
//	  "paging": {"next": "2", "total": 23330},
//	  "linked": {}
//	}
//
// Important gotchas handled here:
//  1. The `description` field is NOT a valid courses.v1 field and causes a 404,
//     so we do NOT request it.
//  2. `startDate` is a Unix epoch-ms NUMBER (or absent), not a string. We decode
//     it as json.RawMessage so we can accept BOTH a number and (defensively) a
//     string, then normalize to an ISO date string for shared.CourseSource.StartDate.
func (s *CourseraSource) fetchFromAPI(ctx context.Context) ([]shared.CourseSource, error) {
	var allCourses []shared.CourseSource

	// Fetch a few pages so we actually get usable volume, not just 20 courses.
	for page := 0; page < 3; page++ {
		offset := page * 20
		url := fmt.Sprintf("%s?limit=20&start=%d&fields=name,slug,startDate,photoUrl,level",
			s.apiURL, offset)

		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
		if err != nil {
			return nil, err
		}

		req.Header.Set("User-Agent", "RojgarSetu/2.0")
		req.Header.Set("Accept", "application/json")

		resp, err := s.client.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("Coursera API returned status: %d for url %s", resp.StatusCode, url)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		var coursData struct {
			Elements []struct {
				Name      string          `json:"name"`
				Slug      string          `json:"slug"`
				Level     string          `json:"level"`
				PhotoURL  string          `json:"photoUrl"`
				StartDate json.RawMessage `json:"startDate"`
			} `json:"elements"`
		}

		if err := json.Unmarshal(body, &coursData); err != nil {
			return nil, err
		}

		if len(coursData.Elements) == 0 {
			break
		}

		for _, c := range coursData.Elements {
			course := shared.CourseSource{
				Source:       "coursera",
				Provider:     "Coursera",
				Title:        shared.CleanString(c.Name),
				URL:          "https://www.coursera.org/learn/" + c.Slug,
				Level:        shared.NormalizeCourseLevel(c.Level),
				ThumbnailURL: c.PhotoURL,
				StartDate:    parseCourseraStartDate(c.StartDate),
				IsFree:       false,
				CreatedAt:    time.Now(),
			}

			if course.Title != "" && shared.IsValidCourse(&course) {
				allCourses = append(allCourses, course)
			}
		}
	}

	log.Info().Int("coursesFromAPI", len(allCourses)).Msg("Coursera API fetch successful")
	return allCourses, nil
}

// parseCourseraStartDate normalizes Coursera's `startDate` into an ISO date
// string pointer for shared.CourseSource.StartDate.
//
// Coursera returns a Unix epoch in MILLISECONDS as a JSON number (e.g.
// 1727880158838), but historically also emitted ISO 8601 strings. We decode the
// raw JSON token and accept both shapes defensively so a future API change does
// not break the whole fetch.
func parseCourseraStartDate(raw json.RawMessage) *string {
	if len(raw) == 0 {
		return nil
	}

	trimmed := strings.TrimSpace(string(raw))

	// Shape 1: numeric epoch milliseconds (the current API shape).
	if trimmed != "null" && trimmed != "" && trimmed[0] != '"' {
		if ms, err := strconv.ParseInt(trimmed, 10, 64); err == nil {
			t := time.UnixMilli(ms)
			iso := t.UTC().Format(time.RFC3339)
			return &iso
		}
		// Not a number; fall through and try string parsing below.
	}

	// Shape 2: string (ISO 8601 or date-only) — defensive for API shape changes.
	var dateStr string
	if err := json.Unmarshal(raw, &dateStr); err != nil {
		return nil
	}
	dateStr = strings.TrimSpace(dateStr)
	if dateStr == "" || dateStr == "null" {
		return nil
	}

	parsed := shared.ParseDateString(dateStr)
	if parsed != nil {
		return parsed
	}
	return nil
}

// fetchFromWebsite fetches from Coursera website
func (s *CourseraSource) fetchFromWebsite(ctx context.Context) ([]shared.CourseSource, error) {
	urls := []string{
		"https://www.coursera.org/browse",
		"https://www.coursera.org/courses",
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

		courses := s.parseHTMLCourses(string(body))
		allCourses = append(allCourses, courses...)
	}

	log.Info().Int("coursesFromWebsite", len(allCourses)).Msg("Coursera website fetch successful")
	return allCourses, nil
}

// parseHTMLCourses parses courses from Coursera HTML
func (s *CourseraSource) parseHTMLCourses(html string) []shared.CourseSource {
	var courses []shared.CourseSource

	// Look for JSON-LD structured data
	pattern := `<script[^>]*type="application/ld\+json"[^>]*>([^<]+)</script>`
	matches := shared.ExtractMatches(html, pattern)

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
				course := shared.CourseSource{
					Source:      "coursera",
					Provider:    "Coursera",
					Title:       shared.CleanString(courseData.Name),
					URL:         courseData.URL,
					Level:       shared.NormalizeCourseLevel(courseData.Level),
					Description: shared.CleanString(courseData.Description),
					IsFree:      false,
					CreatedAt:   time.Now(),
				}
				if shared.IsValidCourse(&course) {
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

