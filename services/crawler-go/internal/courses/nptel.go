package courses

import (
	"github.com/rojgarsetu/crawler/internal/shared"
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
	shared.BaseSource
	client *http.Client
	apiURL string
}

// NewNPTELSource creates a new NPTEL source
func NewNPTELSource() *NPTELSource {
	return &NPTELSource{
		BaseSource: shared.BaseSource{NameStr: "nptel", BaseURL: "https://nptel.ac.in"},
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		// Verified live: the old /api/course.list endpoint returns 404 (gone).
		// The public course-listing HTML page /courses returns 200. NPTEL is an
		// Angular SPA, so the static HTML may not contain every course card, but
		// we extract what anchors/JSON-LD data is present and log honestly.
		apiURL: "https://nptel.ac.in/courses",
	}
}

// Fetch retrieves courses from NPTEL
func (s *NPTELSource) Fetch(ctx context.Context) ([]shared.CourseSource, error) {
	log.Info().Msg("Starting crawl for source: NPTEL")

	var courses []shared.CourseSource

	// Try the HTML course listing first (the old API is dead).
	apiCourses, err := s.fetchFromWebsite(ctx)
	if err != nil {
		log.Error().Err(err).Msg("NPTEL website fetch failed")
		return nil, fmt.Errorf("failed to fetch from NPTEL: %w", err)
	}

	courses = append(courses, apiCourses...)
	log.Info().Int("totalCourses", len(courses)).Msg("NPTEL fetch completed")
	return courses, nil
}

// fetchFromWebsite fetches NPTEL course listings from the HTML pages.
func (s *NPTELSource) fetchFromWebsite(ctx context.Context) ([]shared.CourseSource, error) {
	urls := []string{
		"https://nptel.ac.in/courses",
		"https://nptel.ac.in/courses.html",
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
			log.Warn().Err(err).Str("url", url).Msg("NPTEL request error")
			continue
		}

		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			log.Warn().Int("status", resp.StatusCode).Str("url", url).Msg("NPTEL non-200 response")
			continue
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			continue
		}

		html := string(body)

		// JSON-LD structured data (when present).
		courses := s.parseJSONLDCourses(html)
		allCourses = append(allCourses, courses...)

		// Anchor-based extraction (broad patterns).
		courses = s.parseHTMLCourses(html)
		allCourses = append(allCourses, courses...)

		if len(allCourses) > 0 {
			break
		}
	}

	log.Info().Int("coursesFromWebsite", len(allCourses)).Msg("NPTEL website fetch successful")
	return allCourses, nil
}

// parseJSONLDCourses parses Course JSON-LD blocks from NPTEL HTML.
func (s *NPTELSource) parseJSONLDCourses(html string) []shared.CourseSource {
	var courses []shared.CourseSource

	pattern := `<script[^>]*type="application/ld\+json"[^>]*>([^<]+)</script>`
	matches := shared.ExtractMatches(html, pattern)

	for _, match := range matches {
		if len(match) < 2 {
			continue
		}
		jsonData := match[1]

		var courseData struct {
			Name        string `json:"name"`
			Description string `json:"description"`
			URL         string `json:"url"`
		}
		if err := json.Unmarshal([]byte(jsonData), &courseData); err != nil {
			continue
		}
		if courseData.Name == "" {
			continue
		}

		course := shared.CourseSource{
			Source:      "nptel",
			Provider:    "NPTEL",
			Title:       shared.CleanString(courseData.Name),
			URL:         courseData.URL,
			Description: shared.CleanString(courseData.Description),
			IsFree:      true,
			CreatedAt:   time.Now(),
		}
		if course.URL == "" {
			course.URL = "https://nptel.ac.in/courses"
		}
		if shared.IsValidCourse(&course) {
			courses = append(courses, course)
		}
	}

	return courses
}

// parseHTMLCourses parses courses from NPTEL HTML using broad anchor patterns.
func (s *NPTELSource) parseHTMLCourses(html string) []shared.CourseSource {
	var courses []shared.CourseSource

	patterns := []string{
		`<a[^>]*href="(/course/[^"]*)"[^>]*>.*?<h3[^>]*>([^<]*)</h3>`,
		`<a[^>]*href="(/course/[^"]*)"[^>]*>(?:<[^>]*>)*\s*([^<]{8,120})\s*</a>`,
		`<a[^>]*href="([^"]*course[^"]*)"[^>]*>([^<]{8,120})</a>`,
		`<h3[^>]*>([^<]{8,120})</h3>`,
	}

	for _, pattern := range patterns {
		matches := shared.ExtractMatches(html, pattern)
		for _, match := range matches {
			if len(match) < 2 {
				continue
			}
			// Determine link and title from the match.
			var link, title string
			if len(match) >= 3 {
				link = strings.TrimSpace(match[1])
				title = strings.TrimSpace(match[2])
			} else {
				title = strings.TrimSpace(match[1])
			}

			title = shared.CleanString(title)
			if len(title) < 8 {
				continue
			}

			// Skip navigation / UI text that is not a course.
			if isNonCourseText(title) {
				continue
			}

			course := shared.CourseSource{
				Source:    "nptel",
				Provider:  "NPTEL",
				Title:     title,
				URL:       "https://nptel.ac.in/courses",
				IsFree:    true,
				CreatedAt: time.Now(),
			}
			if link != "" {
				if strings.HasPrefix(link, "http") {
					course.URL = link
				} else if strings.HasPrefix(link, "/") {
					course.URL = "https://nptel.ac.in" + link
				}
			}
			if shared.IsValidCourse(&course) {
				courses = append(courses, course)
			}
		}
	}

	return courses
}

// isNonCourseText filters out navigation/UI strings that are not course titles.
func isNonCourseText(s string) bool {
	lower := strings.ToLower(s)
	skip := []string{
		"courses", "home", "about", "contact", "login", "register",
		"sign in", "sign up", "search", "discipline", "department",
		"nptel", "courses", "lecture", "video", "self-study", "download",
		"certificate", "enrollment", "all courses", "sort by",
	}
	for _, k := range skip {
		if lower == k || strings.Contains(lower, k+" ") {
			// Only skip if it's a short exact-ish match; course titles are long.
			if len(s) < 40 {
				return true
			}
		}
	}
	return false
}

// Name returns the source name
func (s *NPTELSource) Name() string {
	return s.NameStr
}

