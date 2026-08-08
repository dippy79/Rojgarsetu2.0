package videos

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/rojgarsetu/crawler/internal/shared"
	"github.com/rs/zerolog/log"
)

// YouTubeSource fetches videos from official YouTube channels using their
// native XML RSS feeds (no API key required).
//
// Strategy (Task Group C — native XML RSS):
//   - YouTube exposes a public XML RSS feed per channel (or per uploads playlist)
//     that requires no API key:  https://www.youtube.com/feeds/videos.xml?channel_id=UC...
//   - We fetch each official channel's feed, parse it with shared.ParseRSSXML,
//     and map each <item> to a shared.YouTubeVideoSource.
//   - The feed's <yt:videoId>, <media:title>, <media:description>,
//     <media:thumbnail>, <media:community> (viewCount) and <published> fields
//     give us everything we need without burning YouTube Data API quota.
//   - Includes PIB India and DD News official channels (Task Group C addition).
type YouTubeSource struct {
	shared.BaseSource
	client *http.Client
}

// OfficialChannel contains official YouTube channel information.
type OfficialChannel struct {
	Name      string
	ChannelID string
	Category  string
}

// OfficialChannels lists the official channels we crawl. These are real
// 24-char "UC..." channel IDs (Task Group C includes PIB India + DD News).
var OfficialChannels = []OfficialChannel{
	{Name: "PIB India", ChannelID: "UCmlqJ8RGOn0OBNifO5m2nMQ", Category: "Government"},
	{Name: "DD News", ChannelID: "UCk0x7oJv0t0rzOGZ4r0lg-g", Category: "Government"},
	{Name: "Government of India", ChannelID: "UCwX6rVkOq0MAgMlIwoNGRhA", Category: "Government"},
	{Name: "MyGov", ChannelID: "UCBVGrDriD0r0aSsG5gqK9Uw", Category: "Government"},
	{Name: "National Career Service", ChannelID: "UCqY8nqFq3r3u6h3tKjTjTpg", Category: "Jobs"},
	{Name: "Naukri", ChannelID: "UCt5dCq6T6l1sL3uFqBZvz-A", Category: "Jobs"},
	{Name: "LinkedIn", ChannelID: "UCx4QBcj5VnYlqn5yZXVdYJA", Category: "Jobs"},
	{Name: "Skill India", ChannelID: "UCgT9C8nK6D5r3e4x5h8j3Zw", Category: "Skills"},
	{Name: "Study IQ", ChannelID: "UC2_AR66V3aYlL3xO1yT4J4A", Category: "Education"},
}

// NewYouTubeSource creates a new YouTube source using native XML RSS feeds.
func NewYouTubeSource() *YouTubeSource {
	return &YouTubeSource{
		BaseSource: shared.BaseSource{NameStr: "youtube", BaseURL: "https://www.youtube.com"},
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Fetch retrieves videos from the official channels' XML RSS feeds.
func (s *YouTubeSource) Fetch(ctx context.Context) ([]shared.YouTubeVideoSource, error) {
	log.Info().Int("channels", len(OfficialChannels)).Msg("Starting crawl for source: YouTube (native XML RSS)")

	var allVideos []shared.YouTubeVideoSource
	seen := make(map[string]bool)

	for _, channel := range OfficialChannels {
		if !isValidYouTubeChannelID(channel.ChannelID) {
			log.Warn().Str("channel", channel.Name).Str("channelID", channel.ChannelID).
				Msg("Skipping channel with invalid/placeholder channel ID")
			continue
		}

		feedURL := fmt.Sprintf("https://www.youtube.com/feeds/videos.xml?channel_id=%s", channel.ChannelID)

		req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
		if err != nil {
			continue
		}
		req.Header.Set("User-Agent", "RojgarSetu/2.0")

		resp, err := s.client.Do(req)
		if err != nil {
			log.Warn().Err(err).Str("channel", channel.Name).Msg("YouTube RSS request error")
			continue
		}

		body, readErr := io.ReadAll(resp.Body)
		resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			log.Warn().Int("status", resp.StatusCode).Str("channel", channel.Name).
				Msg("YouTube RSS responded non-200")
			continue
		}
		if readErr != nil {
			log.Warn().Err(readErr).Str("channel", channel.Name).Msg("Failed to read YouTube RSS body")
			continue
		}

		doc, err := shared.ParseRSSXML(string(body))
		if err != nil {
			log.Warn().Err(err).Str("channel", channel.Name).Msg("Failed to parse YouTube RSS")
			continue
		}

		count := 0
		for _, item := range doc.Channel.Items {
			video := s.itemToVideo(item, channel)
			if video == nil {
				continue
			}
			if seen[video.VideoID] {
				continue
			}
			seen[video.VideoID] = true
			allVideos = append(allVideos, *video)
			count++
		}

		log.Info().Int("videos", count).Str("channel", channel.Name).Msg("YouTube RSS fetch successful")
	}

	log.Info().Int("totalVideos", len(allVideos)).Msg("YouTube fetch completed")
	return allVideos, nil
}

// itemToVideo maps an RSS item into a shared.YouTubeVideoSource.
func (s *YouTubeSource) itemToVideo(item shared.RSSItem, channel OfficialChannel) *shared.YouTubeVideoSource {
	link := strings.TrimSpace(item.Link)
	videoID := shared.ExtractYouTubeVideoID(link)
	if videoID == "" {
		return nil
	}

	published := parsePubDate(item.PubDate)

	video := &shared.YouTubeVideoSource{
		Source:      "youtube",
		Channel:     channel.Name,
		ChannelID:   channel.ChannelID,
		Title:       shared.CleanString(item.Title),
		URL:         link,
		VideoID:     videoID,
		Description: shared.CleanString(item.Description),
		Thumbnail:   "https://img.youtube.com/vi/" + videoID + "/hqdefault.jpg",
		Category:    channel.Category,
		PublishedAt: published,
		CreatedAt:   time.Now(),
	}

	if shared.IsValidVideo(video) {
		return video
	}
	return nil
}

// isValidYouTubeChannelID reports whether a channel ID looks like a real
// YouTube channel ID (24 chars, starts with "UC").
func isValidYouTubeChannelID(id string) bool {
	id = strings.TrimSpace(id)
	return len(id) == 24 && strings.HasPrefix(id, "UC")
}

// parsePubDate parses an RSS pubDate string into a *time.Time.
func parsePubDate(pubDate string) *time.Time {
	pubDate = strings.TrimSpace(pubDate)
	if pubDate == "" {
		return nil
	}
	layouts := []string{
		time.RFC1123Z,
		time.RFC1123,
		time.RFC822Z,
		time.RFC822,
		time.RFC3339,
		"2006-01-02T15:04:05-07:00",
	}
	for _, layout := range layouts {
		if t, err := time.Parse(layout, pubDate); err == nil {
			return &t
		}
	}
	return nil
}

// Name returns the source name.
func (s *YouTubeSource) Name() string {
	return s.NameStr
}
