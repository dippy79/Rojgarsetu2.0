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

// TutorialsPointSource scrapes courses from TutorialsPoint.
//
// It implements the CourseFetcher interface (Fetch + Name) and writes into the
// shared.CourseSource model, so it is schema-compatible with migration 000012
// (the `courses` table). The store.SaveCourse upsert is idempotent on
// (source, url).
//
// Strategy:
//  1. Fetch the TutorialsPoint course listing page.
//  2. Parse anchors to /tutorials/* and the course catalog section.
//  3. Deduplicate by URL; validate with shared.IsValidCourse.
type TutorialsPointSource struct {
	shared.BaseSource
	client *http.Client
}

// NewTutorialsPointSource creates a new TutorialsPoint source.
func NewTutorialsPointSource() *TutorialsPointSource {
	return &TutorialsPointSource{
		BaseSource: shared.BaseSource{NameStr: "tutorialspoint", BaseURL: "https://www.tutorialspoint.com"},
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Fetch retrieves courses from TutorialsPoint.
func (s *TutorialsPointSource) Fetch(ctx context.Context) ([]shared.CourseSource, error) {
	log.Info().Msg("Starting crawl for source: TutorialsPoint")

	urls := []string{
		"https://www.tutorialspoint.com/courses/",
		"https://www.tutorialspoint.com/tutorials/",
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
			log.Warn().Err(err).Str("url", url).Msg("TutorialsPoint request error")
			continue
		}

		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			log.Warn().Int("status", resp.StatusCode).Str("url", url).Msg("TutorialsPoint non-200 response")
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

		if len(allCourses) > 40 {
			break
		}
	}

	log.Info().Int("totalCourses", len(allCourses)).Msg("TutorialsPoint fetch completed")
	return allCourses, nil
}

// parseCourses extracts courses from TutorialsPoint HTML using broad anchor/H3
// patterns covering both the course catalog and tutorial listing pages.
func (s *TutorialsPointSource) parseCourses(html string) []shared.CourseSource {
	var courses []shared.CourseSource

	patterns := []string{
		`<a[^>]*href="(/tutorials/[^"]*)"[^>]*>(?:<[^>]*>)*([^<]{6,120})</a>`,
		`<a[^>]*href="(/courses/[^"]*)"[^>]*>(?:<[^>]*>)*([^<]{6,120})</a>`,
		`<a[^>]*href="([^"]*tutorialspoint\.com/tutorials/[^"]*)"[^>]*>([^<]{6,120})</a>`,
	}

	for _, pattern := range patterns {
		matches := shared.ExtractMatches(html, pattern)
		for _, match := range matches {
			if len(match) < 3 {
				continue
			}
			link := strings.TrimSpace(match[1])
			title := shared.CleanString(match[2])

			if len(title) < 6 || isTPNonCourseText(title) {
				continue
			}

			url := link
			if strings.HasPrefix(link, "/") {
				url = "https://www.tutorialspoint.com" + link
			}

			course := shared.CourseSource{
				Source:    "tutorialspoint",
				Provider:  "TutorialsPoint",
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

// isTPNonCourseText filters out navigation/UI strings that are not courses.
func isTPNonCourseText(s string) bool {
	lower := strings.ToLower(s)
	skip := []string{
		"home", "about", "contact", "login", "register", "sign in",
		"sign up", "search", "courses", "all courses", "tutorials library",
		"certification", "training", "print page", "top", "career",
		"library", "category", "browse", "free tutorials", "jobs",
	}
	for _, k := range skip {
		if lower == k {
			return true
		}
	}
	return false
}

// Name returns the source name.
func (s *TutorialsPointSource) Name() string {
	return s.NameStr
}
