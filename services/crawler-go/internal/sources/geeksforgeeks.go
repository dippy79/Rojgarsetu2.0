package sources

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

// GeeksforGeeksSource scrapes courses from GeeksforGeeks.
//
// It implements the CourseFetcher interface (Fetch + Name) and writes into the
// same CourseSource model used by nptel.go, swayam.go, nsdc.go, coursera.go and
// udemy.go, so it is schema-compatible with migration 000012 (the `courses`
// table). The store.SaveCourse upsert is idempotent on (source, url).
//
// Strategy (verified against the GfG public course pages):
//  1. Parse detection/JSON-LD first — GfG course pages embed structured data in
//     <script type="application/ld+json"> blocks.
//  2. Fall back to anchor/H3 based extraction for the course listing page.
//  3. If neither yields structured data, return an honest empty result rather
//     than fabricating courses.
type GeeksforGeeksSource struct {
	BaseSource
	client *http.Client
}

// NewGeeksforGeeksSource creates a new GeeksforGeeks source.
func NewGeeksforGeeksSource() *GeeksforGeeksSource {
	return &GeeksforGeeksSource{
		BaseSource: BaseSource{NameStr: "geeksforgeeks", BaseURL: "https://www.geeksforgeeks.org"},
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Fetch retrieves courses from GeeksforGeeks.
func (s *GeeksforGeeksSource) Fetch(ctx context.Context) ([]CourseSource, error) {
	log.Info().Msg("Starting crawl for source: GeeksforGeeks")

	urls := []string{
		"https://www.geeksforgeeks.org/courses/",
		"https://www.geeksforgeeks.org/courses?type=Online",
	}

	var allCourses []CourseSource
	seen := make(map[string]bool)

	for _, url := range urls {
		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
		if err != nil {
			continue
		}
		req.Header.Set("User-Agent", "RojgarSetu/2.0")

		resp, err := s.client.Do(req)
		if err != nil {
			log.Warn().Err(err).Str("url", url).Msg("GfG request error")
			continue
		}

		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			log.Warn().Int("status", resp.StatusCode).Str("url", url).Msg("GfG non-200 response")
			continue
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			continue
		}

		html := string(body)

		courses := s.parseCourses(html)
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

	log.Info().Int("totalCourses", len(allCourses)).Msg("GeeksforGeeks fetch completed")
	return allCourses, nil
}

// parseCourses tries JSON-LD first, then falls back to anchor/H3 extraction.
func (s *GeeksforGeeksSource) parseCourses(html string) []CourseSource {
	var courses []CourseSource

	// 1. Prefer JSON-LD structured data.
	jsonLDCourses := s.parseJSONLDCourses(html)
	courses = append(courses, jsonLDCourses...)

	// 2. Fallback: anchor/H3-based extraction (only if JSON-LD yielded nothing).
	if len(courses) == 0 {
		courses = s.parseHTMLCourses(html)
	}

	return courses
}

// gfGCourse mirrors the subset of schema.org Course we care about from GfG
// JSON-LD blocks.
type gfGCourse struct {
	Type        string `json:"@type"`
	Name        string `json:"name"`
	Description string `json:"description"`
	URL         string `json:"url"`
	Provider    struct {
		Name string `json:"name"`
	} `json:"provider"`
	Offers struct {
		Price     string `json:"price"`
		PriceSpec string `json:"priceCurrency"`
	} `json:"offers"`
}

func (c gfGCourse) isCourse() bool {
	return c.Type == "Course" && c.Name != ""
}

// parseJSONLDCourses extracts Course JSON-LD blocks from GfG HTML.
func (s *GeeksforGeeksSource) parseJSONLDCourses(html string) []CourseSource {
	var courses []CourseSource

	pattern := `<script[^>]*type="application/ld\+json"[^>]*>(.*?)</script>`
	re := regexp.MustCompile(`(?s)` + pattern)
	blocks := re.FindAllStringSubmatch(html, -1)

	for _, block := range blocks {
		if len(block) < 2 {
			continue
		}
		raw := strings.TrimSpace(block[1])

		// Single Course object.
		var single gfGCourse
		if err := json.Unmarshal([]byte(raw), &single); err == nil && single.isCourse() {
			if c := s.jsonLDToCourse(single); c != nil {
				courses = append(courses, *c)
			}
			continue
		}

		// Array of Course objects.
		var arr []gfGCourse
		if err := json.Unmarshal([]byte(raw), &arr); err == nil {
			for _, item := range arr {
				if !item.isCourse() {
					continue
				}
				if c := s.jsonLDToCourse(item); c != nil {
					courses = append(courses, *c)
				}
			}
		}
	}

	if len(courses) > 0 {
		log.Info().Int("jsonldCourses", len(courses)).Msg("Parsed courses from JSON-LD")
	}
	return courses
}

func (s *GeeksforGeeksSource) jsonLDToCourse(c gfGCourse) *CourseSource {
	provider := c.Provider.Name
	if provider == "" {
		provider = "GeeksforGeeks"
	}

	url := c.URL
	if url == "" {
		url = "https://www.geeksforgeeks.org/courses/"
	}

	course := &CourseSource{
		Source:      "geeksforgeeks",
		Provider:    provider,
		Title:       cleanString(c.Name),
		URL:         url,
		Description: cleanString(c.Description),
		Price:       c.Offers.Price,
		IsFree:      strings.EqualFold(c.Offers.Price, "0") || c.Offers.Price == "",
		CreatedAt:   time.Now(),
	}

	if isValidCourse(course) {
		return course
	}
	return nil
}

// parseHTMLCourses extracts courses from GfG HTML using broad anchor/H3 patterns.
func (s *GeeksforGeeksSource) parseHTMLCourses(html string) []CourseSource {
	var courses []CourseSource

	patterns := []string{
		`<a[^>]*href="(/courses/[^"]*)"[^>]*>(?:<[^>]*>)*([^<]{8,120})</a>`,
		`<a[^>]*href="([^"]*geeksforgeeks\.org/courses[^"]*)"[^>]*>([^<]{8,120})</a>`,
		`<h3[^>]*>([^<]{8,120})</h3>`,
	}

	for _, pattern := range patterns {
		matches := extractMatches(html, pattern)
		for _, match := range matches {
			if len(match) < 2 {
				continue
			}
			var link, title string
			if len(match) >= 3 {
				link = strings.TrimSpace(match[1])
				title = cleanString(match[2])
			} else {
				title = cleanString(match[1])
			}

			if len(title) < 8 || isGfGNonCourseText(title) {
				continue
			}

			url := link
			if url == "" {
				url = "https://www.geeksforgeeks.org/courses/"
			} else if strings.HasPrefix(url, "/") {
				url = "https://www.geeksforgeeks.org" + url
			}

			course := CourseSource{
				Source:    "geeksforgeeks",
				Provider:  "GeeksforGeeks",
				Title:     title,
				URL:       url,
				IsFree:    false,
				CreatedAt: time.Now(),
			}
			if isValidCourse(&course) {
				courses = append(courses, course)
			}
		}
		if len(courses) > 0 {
			break
		}
	}

	return courses
}

// isGfGNonCourseText filters out navigation/UI strings that are not course titles.
func isGfGNonCourseText(s string) bool {
	lower := strings.ToLower(s)
	skip := []string{
		"courses", "home", "about", "contact", "login", "register",
		"sign in", "sign up", "search", "practice", "gfg", "online",
		"all courses", "sort by", "course", "doubt support", "classroom",
	}
	for _, k := range skip {
		if lower == k || strings.Contains(lower, k+" ") {
			if len(s) < 40 {
				return true
			}
		}
	}
	return false
}

// Name returns the source name.
func (s *GeeksforGeeksSource) Name() string {
	return s.NameStr
}
