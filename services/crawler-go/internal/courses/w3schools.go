package courses

import (
	"context"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/rojgarsetu/crawler/internal/shared"
	"github.com/rs/zerolog/log"
)

// W3SchoolsSource scrapes courses from W3Schools.
//
// It implements the shared.CourseFetcher interface (Fetch + Name) and writes
// into the shared.CourseSource model, so it is schema-compatible with migration
// 000012 (the `courses` table). The store.SaveCourse upsert is idempotent on
// (source, url).
//
// Strategy:
//  1. Fetch the W3Schools tutorial index page.
//  2. Parse anchors to tutorial paths (e.g. /html/, /css/, /python/) and
//     the course/tutorial catalog sections.
//  3. Deduplicate by URL; validate with shared.IsValidCourse.
type W3SchoolsSource struct {
	shared.BaseSource
	client *http.Client
}

// NewW3SchoolsSource creates a new W3Schools source.
func NewW3SchoolsSource() *W3SchoolsSource {
	return &W3SchoolsSource{
		BaseSource: shared.BaseSource{NameStr: "w3schools", BaseURL: "https://www.w3schools.com"},
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Fetch retrieves courses from W3Schools.
func (s *W3SchoolsSource) Fetch(ctx context.Context) ([]shared.CourseSource, error) {
	log.Info().Msg("Starting crawl for source: W3Schools")

	urls := []string{
		"https://www.w3schools.com/",
		"https://www.w3schools.com/tutorials/",
	}

	var allCourses []shared.CourseSource
	seen := make(map[string]bool)

	for _, url := range urls {
		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
		if err != nil {
			continue
		}
		req.Header.Set("User-Agent", "RojgarSetu/2.0")

		resp, err := s.client.Do(req)
		if err != nil {
			log.Warn().Err(err).Str("url", url).Msg("W3Schools request error")
			continue
		}

		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			log.Warn().Int("status", resp.StatusCode).Str("url", url).Msg("W3Schools non-200 response")
			continue
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			continue
		}

		courses := s.parseCourses(string(body))
		for _, c := range courses {
			key := c.URL
			if key == "" {
				key = c.Title
			}
			if !seen[key] {
				seen[key] = true
				allCourses = append(allCourses, c)
			}
		}

		if len(allCourses) > 0 {
			break
		}
	}

	log.Info().Int("totalCourses", len(allCourses)).Msg("W3Schools fetch completed")
	return allCourses, nil
}

// parseCourses extracts tutorial/course cards from W3Schools HTML.
func (s *W3SchoolsSource) parseCourses(html string) []shared.CourseSource {
	var courses []shared.CourseSource

	patterns := []string{
		`<a[^>]*href="(/(?:html|css|js|python|java|sql|c|php|bootstrap|react|nodejs|w3css|jquery|xml|django|angular|typescript|r|go|kotlin|dart|tailwind|aws|azure|sass|git|golang)/)[^"]*"[^>]*>(?:<[^>]*>)*([^<]{2,80})</a>`,
		`<a[^>]*href="([^"]*w3schools\.com/[^"]*)"[^>]*>(?:<[^>]*>)*([^<]{2,80})</a>`,
	}

	for _, pattern := range patterns {
		matches := shared.ExtractMatches(html, pattern)
		for _, match := range matches {
			if len(match) < 3 {
				continue
			}
			link := strings.TrimSpace(match[1])
			title := shared.CleanString(match[2])

			if len(title) < 2 || isW3NonCourseText(title) {
				continue
			}

			url := link
			if strings.HasPrefix(link, "/") {
				url = "https://www.w3schools.com" + link
			}

			course := shared.CourseSource{
				Source:    "w3schools",
				Provider:  "W3Schools",
				Title:     title,
				URL:       url,
				IsFree:    true,
				Mode:      "online",
				CreatedAt: time.Now(),
			}
			if shared.IsValidCourse(&course) {
				courses = append(courses, course)
			}
		}
		if len(courses) > 0 {
			break
		}
	}

	return courses
}

// isW3NonCourseText filters out navigation/UI strings that are not courses.
func isW3NonCourseText(s string) bool {
	lower := strings.ToLower(s)
	skip := []string{
		"home", "about", "contact", "login", "register", "sign in",
		"sign up", "search", "w3schools", "spaces", "certificates",
		"tutorials", "references", "exercises", "learn", "get certified",
		"how to", "icons", "colors", "game", "browser", "forum", "jobs",
	}
	for _, k := range skip {
		if strings.TrimSpace(lower) == k {
			return true
		}
	}
	return false
}

// Name returns the source name.
func (s *W3SchoolsSource) Name() string {
	return s.NameStr
}
