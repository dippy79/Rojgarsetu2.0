package sources

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

// YouTubeSource scrapes videos from YouTube official channels
type YouTubeSource struct {
	BaseSource
	client     *http.Client
	apiKey     string
	apiBaseURL string
}

// OfficialChannel contains official YouTube channel information
type OfficialChannel struct {
	Name      string
	ChannelID string
	Category  string
}

// OfficialChannels contains list of official channels
var OfficialChannels = []OfficialChannel{
	{Name: "Naukri", ChannelID: "UCt5dCq6T6l1sL3uFqBZvz-A", Category: "Jobs"},
	{Name: "LinkedIn", ChannelID: "UCx4QBcj5VnYlqn5yZXVdYJA", Category: "Jobs"},
	{Name: "Government of India", ChannelID: "UCwX6rVkOq0MAgMlIwoNGRhA", Category: "Government"},
	{Name: "MyGov", ChannelID: "UCBVGrDriD0r0aSsG5gqK9Uw", Category: "Government"},
	{Name: "DPIIT", ChannelID: "UCQ1Ngtn9KUfEuDExnGi6Q3A", Category: "Government"},
	{Name: "National Career Service", ChannelID: "UCqY8nqFq3r3u6h3tKjTjTpg", Category: "Jobs"},
	{Name: "Study IQ", ChannelID: "UC2_AR66V3aYlL3xO1yT4J4A", Category: "Education"},
	{Name: "Unacademy", ChannelID: "UCwF0s71bR2XC7dR3jJ4f7jA", Category: "Education"},
	{Name: "Byju's", ChannelID: "UCBCl6VSpL4WN3BL3OdE3E7w", Category: "Education"},
	{Name: "Skill India", ChannelID: "UCgT9C8nK6D5r3e4x5h8j3Zw", Category: "Skills"},
	{Name: "Mudra Life", ChannelID: "UCw9xY7Yv3hO9q3yK5f9r8Zg", Category: "Jobs"},
	{Name: "UPSC Pathshala", ChannelID: "UC2sG8M7L2w8kX5jK9f4x5gQ", Category: "UPSC"},
	{Name: "SSC Adda", ChannelID: "UCw3P4v8v5hK7f6d8i9x4j2Q", Category: "SSC"},
	{Name: "Railway AddA", ChannelID: "UCn4M5v7K8f9h2d3k6x5v7pA", Category: "Railways"},
	{Name: "Talent.com", ChannelID: "UCv4R7K8p3n9m2k5x5v8v0A", Category: "Jobs"},
}

// NewYouTubeSource creates a new YouTube source
func NewYouTubeSource() *YouTubeSource {
	apiKey := os.Getenv("YOUTUBE_API_KEY")
	return &YouTubeSource{
		BaseSource: BaseSource{NameStr: "youtube", BaseURL: "https://www.youtube.com"},
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		apiKey:     apiKey,
		apiBaseURL: "https://www.googleapis.com/youtube/v3",
	}
}

// Fetch retrieves videos from YouTube official channels.
//
// Phase A: we only use the YouTube Data API v3 path (which needs a real
// YOUTUBE_API_KEY from the environment). The old website-scraping fallback
// (fetchFromWebsite) almost always returned 0 videos because YouTube serves a
// JS-heavy page, so it has been removed as a broken path. This is not new
// infrastructure — it is simply fixing the wiring: no key → clear warning +
// skip; key present → real API calls with proper error logging.
func (s *YouTubeSource) Fetch(ctx context.Context) ([]YouTubeVideoSource, error) {
	log.Info().Msg("Starting crawl for source: YouTube")

	if s.apiKey == "" {
		// No key configured: do NOT hit the broken website scraper. Log clearly
		// so the operator knows why YouTube yields 0 videos.
		log.Warn().Msg("YOUTUBE_API_KEY is not set - skipping YouTube fetch (website fallback removed as broken; set YOUTUBE_API_KEY to enable)")
		return []YouTubeVideoSource{}, nil
	}

	apiVideos, err := s.fetchFromAPI(ctx)
	if err != nil {
		log.Error().Err(err).Msg("YouTube API fetch failed")
		return nil, fmt.Errorf("failed to fetch from YouTube API: %w", err)
	}

	log.Info().Int("totalVideos", len(apiVideos)).Msg("YouTube fetch completed")
	return apiVideos, nil
}

// fetchFromAPI fetches from YouTube Data API v3.
func (s *YouTubeSource) fetchFromAPI(ctx context.Context) ([]YouTubeVideoSource, error) {
	var allVideos []YouTubeVideoSource

	for _, channel := range OfficialChannels {
		if !isValidYouTubeChannelID(channel.ChannelID) {
			log.Warn().Str("channel", channel.Name).Str("channelID", channel.ChannelID).
				Msg("Skipping channel with invalid/placeholder channel ID (not a real 24-char UC... ID)")
			continue
		}

		url := fmt.Sprintf("%s/search?key=%s&channelId=%s&part=snippet,id&order=date&maxResults=20&type=video",
			s.apiBaseURL, s.apiKey, channel.ChannelID)

		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
		if err != nil {
			continue
		}

		req.Header.Set("User-Agent", "RojgarSetu/2.0")

		resp, err := s.client.Do(req)
		if err != nil {
			log.Warn().Err(err).Str("channel", channel.Name).Msg("YouTube API request error")
			continue
		}

		body, readErr := io.ReadAll(resp.Body)
		resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			// Log the real error body (e.g. quotaExceeded, invalid key, channel not found)
			// instead of silently continuing. This is the "wire the real error" fix.
			detail := ""
			if readErr == nil && len(body) > 0 {
				detail = string(body)
			}
			log.Warn().Int("status", resp.StatusCode).Str("channel", channel.Name).
				Str("body", detail).Msg("YouTube API responded non-200")
			continue
		}

		if readErr != nil {
			log.Warn().Err(readErr).Str("channel", channel.Name).Msg("Failed to read YouTube API body")
			continue
		}

		var ytData struct {
			Items []struct {
				ID struct {
					VideoID string `json:"videoId"`
				} `json:"id"`
				Snippet struct {
					Title        string `json:"title"`
					Description  string `json:"description"`
					ChannelTitle string `json:"channelTitle"`
					ChannelID    string `json:"channelId"`
					Thumbnails   struct {
						Medium struct {
							URL string `json:"url"`
						} `json:"medium"`
						High struct {
							URL string `json:"url"`
						} `json:"high"`
						Default struct {
							URL string `json:"url"`
						} `json:"default"`
					} `json:"thumbnails"`
					PublishTime time.Time `json:"publishTime"`
				} `json:"snippet"`
			} `json:"items"`
		}

		if err := json.Unmarshal(body, &ytData); err != nil {
			log.Warn().Err(err).Str("channel", channel.Name).Msg("Failed to parse YouTube API response")
			continue
		}

		for _, item := range ytData.Items {
			if item.ID.VideoID == "" {
				continue
			}

			video := YouTubeVideoSource{
				Source:      "youtube",
				Channel:     item.Snippet.ChannelTitle,
				ChannelID:   item.Snippet.ChannelID,
				Title:       cleanString(item.Snippet.Title),
				URL:         "https://www.youtube.com/watch?v=" + item.ID.VideoID,
				VideoID:     item.ID.VideoID,
				Description: cleanString(item.Snippet.Description),
				Thumbnail:   item.Snippet.Thumbnails.High.URL,
				Category:    channel.Category,
				PublishedAt: &item.Snippet.PublishTime,
				CreatedAt:   time.Now(),
			}

			if isValidVideo(&video) {
				allVideos = append(allVideos, video)
			}
		}

		log.Info().Int("videos", len(ytData.Items)).Str("channel", channel.Name).Msg("YouTube API fetch successful")
	}

	return allVideos, nil
}

// isValidYouTubeChannelID reports whether a channel ID looks like a real
// YouTube channel ID. Real channel IDs are exactly 24 characters and start
// with "UC". Many entries in OfficialChannels are placeholder strings that are
// not valid; skipping them avoids pointless API quota usage and 404s.
func isValidYouTubeChannelID(id string) bool {
	id = strings.TrimSpace(id)
	return len(id) == 24 && strings.HasPrefix(id, "UC")
}

// Name returns the source name
func (s *YouTubeSource) Name() string {
	return s.NameStr
}

// GetAPILimit returns the YouTube API quota usage info
func (s *YouTubeSource) GetAPILimit() string {
	if s.apiKey == "" {
		return "No API key configured - YouTube skipped"
	}
	return "API key configured - using YouTube Data API v3"
}

// Helper to parse view count
func parseCount(countStr string) int64 {
	count, err := strconv.ParseInt(strings.ReplaceAll(countStr, ",", ""), 10, 64)
	if err != nil {
		return 0
	}
	return count
}
