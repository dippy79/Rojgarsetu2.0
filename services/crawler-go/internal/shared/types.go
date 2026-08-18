package shared

import (
	"context"
	"encoding/xml"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// ============================================
// DATA STRUCTURES
// ============================================

// GovJobSource represents a government job from various sources
type GovJobSource struct {
	Source           string         `json:"source"`
	Title            string         `json:"title"`
	Department       string         `json:"department"`
	Location         string         `json:"location"`
	ApplyURL         string         `json:"apply_url"`
	LastDate         *string        `json:"last_date,omitempty"`
	Eligibility      string         `json:"eligibility"`
	VacancyCount     *int           `json:"vacancy_count,omitempty"`
	Salary           string         `json:"salary"`
	ExamDate         *string        `json:"exam_date,omitempty"`
	NotificationURL  string         `json:"notification_url,omitempty"`
	JobRegion        string         `json:"job_region,omitempty"`
	IsVerified       bool           `json:"is_verified"`
	VerificationMeta map[string]any `json:"verification_meta,omitempty"`
	ScamScore        float64        `json:"scam_score"`
	CanonicalURL     string         `json:"canonical_url,omitempty"`
	CreatedAt        time.Time      `json:"created_at"`
}

// PrivJobSource represents a private job from various sources
type PrivJobSource struct {
	Source           string         `json:"source"`
	Company          string         `json:"company"`
	Title            string         `json:"title"`
	Location         string         `json:"location"`
	URL              string         `json:"url"`
	Salary           string         `json:"salary"`
	Experience       string         `json:"experience"`
	JobType          string         `json:"job_type"`
	Skills           []string       `json:"skills"`
	Description      string         `json:"description"`
	PostedAt         *time.Time     `json:"posted_at,omitempty"`
	JobRegion        string         `json:"job_region,omitempty"`
	IsVerified       bool           `json:"is_verified"`
	VerificationMeta map[string]any `json:"verification_meta,omitempty"`
	ScamScore        float64        `json:"scam_score"`
	CanonicalURL     string         `json:"canonical_url,omitempty"`
	CreatedAt        time.Time      `json:"created_at"`
}

// CourseSource represents a course from various providers
type CourseSource struct {
	Source          string    `json:"source"`
	Provider        string    `json:"provider"`
	Title           string    `json:"title"`
	URL             string    `json:"url"`
	Duration        string    `json:"duration"`
	Mode            string    `json:"mode"`
	Level           string    `json:"level"`
	Skills          []string  `json:"skills"`
	Description     string    `json:"description"`
	ThumbnailURL    string    `json:"thumbnail_url"`
	Price           string    `json:"price"`
	IsFree          bool      `json:"is_free"`
	StartDate       *string   `json:"start_date,omitempty"`
	EndDate         *string   `json:"end_date,omitempty"`
	EnrollmentCount *int      `json:"enrollment_count,omitempty"`
	Rating          *float64  `json:"rating,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
}

// YouTubeVideoSource represents a YouTube video from official channels
type YouTubeVideoSource struct {
	Source      string     `json:"source"`
	Channel     string     `json:"channel"`
	ChannelID   string     `json:"channel_id"`
	Title       string     `json:"title"`
	URL         string     `json:"url"`
	Thumbnail   string     `json:"thumbnail"`
	Description string     `json:"description"`
	VideoID     string     `json:"video_id"`
	PublishedAt *time.Time `json:"published_at,omitempty"`
	Duration    string     `json:"duration"`
	ViewCount   *int64     `json:"view_count,omitempty"`
	LikeCount   *int64     `json:"like_count,omitempty"`
	Category    string     `json:"category"`
	CreatedAt   time.Time  `json:"created_at"`
}

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

// ParseRSSXML parses RSS XML content
func ParseRSSXML(xmlContent string) (*RSSDocument, error) {
	var doc RSSDocument
	err := xml.Unmarshal([]byte(xmlContent), &doc)
	if err != nil {
		return nil, err
	}
	return &doc, nil
}

// CleanString removes extra whitespace and trims
func CleanString(s string) string {
	s = strings.TrimSpace(s)
	s = strings.Join(strings.Fields(s), " ")
	return s
}

// ExtractURL extracts URL from href or link
func ExtractURL(link string) string {
	link = strings.TrimSpace(link)
	if strings.HasPrefix(link, "http") {
		return link
	}
	return ""
}

// ExtractField extracts a field value from description text
func ExtractField(text, fieldName string) string {
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

// ParseDateString parses various date formats
func ParseDateString(dateStr string) *string {
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

	// Return original if parsing fails
	return &dateStr
}

// IsValidJob checks if job has minimum required fields
func IsValidJob(job *GovJobSource) bool {
	if job == nil {
		return false
	}
	if job.Title == "" || len(job.Title) < 3 {
		return false
	}
	if job.Source == "" {
		return false
	}
	// Skip if title looks like an error or junk
	skipPatterns := []string{"error", "not found", "404", "maintenance"}
	for _, pattern := range skipPatterns {
		if strings.Contains(strings.ToLower(job.Title), pattern) {
			return false
		}
	}
	return true
}

// IsValidPrivJob checks if private job has minimum required fields
func IsValidPrivJob(job *PrivJobSource) bool {
	if job == nil {
		return false
	}
	if job.Title == "" || len(job.Title) < 3 {
		return false
	}
	if job.Company == "" {
		return false
	}
	if job.Source == "" {
		return false
	}
	return true
}

// IsValidCourse checks if course has minimum required fields
func IsValidCourse(course *CourseSource) bool {
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

// IsValidVideo checks if video has minimum required fields
func IsValidVideo(video *YouTubeVideoSource) bool {
	if video == nil {
		return false
	}
	if video.Title == "" || len(video.Title) < 3 {
		return false
	}
	if video.Channel == "" {
		return false
	}
	if video.VideoID == "" {
		return false
	}
	return true
}

// ExtractMatches extracts matches using regex pattern
func ExtractMatches(text, pattern string) [][]string {
	re := regexp.MustCompile(pattern)
	return re.FindAllStringSubmatch(text, -1)
}

// ExtractYouTubeVideoID extracts video ID from YouTube URL
func ExtractYouTubeVideoID(url string) string {
	patterns := []string{
		`(?:youtube\.com/watch\?v=|youtu\.be/|youtube\.com/embed/)([a-zA-Z0-9_-]{11})`,
		`^([a-zA-Z0-9_-]{11})$`,
	}
	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(url)
		if len(matches) > 1 {
			return matches[1]
		}
	}
	return ""
}

// ParseYouTubeDuration parses YouTube ISO 8601 duration
func ParseYouTubeDuration(duration string) string {
	// P1H2M10S -> 1h 2m 10s
	re := regexp.MustCompile(`PT(?:(\d+)H)?(?:(\d+)M)?(?:(\d+)S)?`)
	matches := re.FindStringSubmatch(duration)
	if len(matches) == 0 {
		return duration
	}

	var parts []string
	if matches[1] != "" {
		parts = append(parts, matches[1]+"h")
	}
	if matches[2] != "" {
		parts = append(parts, matches[2]+"m")
	}
	if matches[3] != "" {
		parts = append(parts, matches[3]+"s")
	}

	if len(parts) == 0 {
		return duration
	}
	return strings.Join(parts, " ")
}

// ParseViewCount parses view count string
func ParseViewCount(countStr string) *int64 {
	countStr = strings.ReplaceAll(countStr, ",", "")
	countStr = strings.TrimSpace(countStr)
	count, err := strconv.ParseInt(countStr, 10, 64)
	if err != nil {
		return nil
	}
	return &count
}

// NormalizeJobType normalizes job type string
func NormalizeJobType(jobType string) string {
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

// NormalizeCourseMode normalizes course mode string
func NormalizeCourseMode(mode string) string {
	mode = strings.ToLower(strings.TrimSpace(mode))

	if strings.Contains(mode, "online") {
		return "online"
	}
	if strings.Contains(mode, "offline") || strings.Contains(mode, "classroom") {
		return "offline"
	}
	if strings.Contains(mode, "hybrid") {
		return "hybrid"
	}

	return mode
}

// NormalizeCourseLevel normalizes course level string
func NormalizeCourseLevel(level string) string {
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

// ============================================
// FETCHER INTERFACES
// ============================================

// GovJobFetcher interface for government job sources
type GovJobFetcher interface {
	Fetch(ctx context.Context) ([]GovJobSource, error)
	Name() string
}

// PrivJobFetcher interface for private job sources
type PrivJobFetcher interface {
	Fetch(ctx context.Context) ([]PrivJobSource, error)
	Name() string
}

// CourseFetcher interface for course sources
type CourseFetcher interface {
	Fetch(ctx context.Context) ([]CourseSource, error)
	Name() string
}

// VideoFetcher interface for YouTube video sources
type VideoFetcher interface {
	Fetch(ctx context.Context) ([]YouTubeVideoSource, error)
	Name() string
}
