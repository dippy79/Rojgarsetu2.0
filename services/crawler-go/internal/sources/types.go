package sources

import (
	"encoding/xml"
	"net/url"
	"regexp"
	"strings"
	"time"
	"github.com/rojgarsetu/crawler/internal/shared"
)

// ============================================
// DATA STRUCTURES
// ============================================

// Job struct used for intermediate processing
type Job struct {
	Title             string
	CompanyOrDept     string
	Location          string
	QualificationReq  string
	SalaryOrPayScale  string
	ApplyURL          string
	SourceAttribution string
	HashChecksum      string
	Category          string // CENTRAL | STATE
	StateName         string // e.g., "Punjab", "Bihar", "ALL_INDIA"
}

// Re-export common types from shared
type GovJobSource = shared.GovJobSource
type PrivJobSource = shared.PrivJobSource
type CourseSource = shared.CourseSource
type YouTubeVideoSource = shared.YouTubeVideoSource
type GovFormSource = shared.GovFormSource
type GovtJobSource = shared.GovJobSource

// ============================================
// RSS STRUCTURES
// ============================================

// RSSDocument represents an RSS feed
type RSSDocument struct {
	XMLName xml.Name   `xml:"rss"`
	Channel RSSChannel `xml:"channel"`
}

// RSSChannel represents an RSS channel
type RSSChannel struct {
	Title       string    `xml:"title"`
	Description string    `xml:"description"`
	Link        string    `xml:"link"`
	Items       []RSSItem `xml:"item"`
}

// RSSItem represents an RSS item
type RSSItem struct {
	Title       string `xml:"title"`
	Description string `xml:"description"`
	Link        string `xml:"link"`
	PubDate     string `xml:"pubDate"`
	Category    string `xml:"category"`
	Author      string `xml:"author"`
	GUID        string `xml:"guid"`
}

// ============================================
// HELPER FUNCTIONS
// ============================================

func parseRSSXML(xmlContent string) (*RSSDocument, error) {
	var doc RSSDocument
	err := xml.Unmarshal([]byte(xmlContent), &doc)
	if err != nil {
		return nil, err
	}
	return &doc, nil
}

func cleanString(s string) string {
	s = strings.TrimSpace(s)
	s = strings.Join(strings.Fields(s), " ")
	return s
}

func extractURL(link string) string {
	link = strings.TrimSpace(link)
	if strings.HasPrefix(link, "http") {
		return link
	}
	return ""
}

func extractField(text, fieldName string) string {
	patterns := []string{
		fieldName + `([^<\n]+)`,
		fieldName + `:\s*([^<\n]+)`,
	}
	for _, pattern := range patterns {
		re := regexp.MustCompile(`(?i)` + pattern)
		matches := re.FindStringSubmatch(text)
		if len(matches) > 1 {
			return strings.TrimSpace(matches[1])
		}
	}
	return ""
}

func parseDateString(dateStr string) *string {
	if dateStr == "" {
		return nil
	}

	dateStr = strings.TrimSpace(dateStr)
	formats := []string{
		"2006-01-02",
		"02/01/2006",
		"02-01-2006",
		"Jan 02, 2006",
		"02 January 2006",
		"02 Jan 2006",
		"2006/01/02",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, dateStr); err == nil {
			result := t.Format("2006-01-02")
			return &result
		}
	}

	return &dateStr
}

func isValidJob(job *Job) bool {
	if job == nil {
		return false
	}
	if job.Title == "" || len(job.Title) < 3 {
		return false
	}
	if job.ApplyURL == "" {
		return false
	}
	return true
}

func isValidGovtJob(job *GovJobSource) bool {
	if job == nil {
		return false
	}
	if job.Title == "" || len(job.Title) < 3 {
		return false
	}
	return true
}

func isValidPrivJob(job *PrivJobSource) bool {
	if job == nil {
		return false
	}
	if job.Title == "" || len(job.Title) < 3 {
		return false
	}
	if job.Company == "" {
		return false
	}
	return true
}

func isValidCourse(course *CourseSource) bool {
	if course == nil {
		return false
	}
	if course.Title == "" || len(course.Title) < 3 {
		return false
	}
	if course.Provider == "" {
		return false
	}
	if course.URL == "" {
		return false
	}
	return true
}

func normalizeJobType(jobType string) string {
	jt := strings.ToLower(strings.TrimSpace(jobType))

	typeMappings := map[string]string{
		"full time":   "full-time",
		"fulltime":    "full-time",
		"part time":   "part-time",
		"parttime":    "part-time",
		"contract":    "contract",
		"contractual": "contract",
		"internship":  "internship",
		"intern":      "internship",
		"temporary":   "temporary",
		"permanent":   "permanent",
		"remote":      "remote",
		"hybrid":      "hybrid",
	}

	for k, v := range typeMappings {
		if strings.Contains(jt, k) {
			return v
		}
	}

	return jobType
}

func normalizeCourseLevel(level string) string {
	level = strings.ToLower(strings.TrimSpace(level))

	if strings.Contains(level, "beginner") || strings.Contains(level, "basic") {
		return "beginner"
	}
	if strings.Contains(level, "intermediate") {
		return "intermediate"
	}
	if strings.Contains(level, "advanced") || strings.Contains(level, "expert") {
		return "advanced"
	}

	return level
}

func extractMatches(text, pattern string) [][]string {
	re := regexp.MustCompile(pattern)
	return re.FindAllStringSubmatch(text, -1)
}

func extractDomainName(rawURL string) string {
	u, err := url.Parse(rawURL)
	if err != nil {
		return rawURL
	}
	return u.Hostname()
}
