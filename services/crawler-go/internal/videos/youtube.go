package videos

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/rojgarsetu/crawler/internal/shared"
	"github.com/rs/zerolog/log"
)

type YouTubeSource struct {
	shared.BaseSource
	client *http.Client
}

type YouTubeSearchResponse struct {
	Items []struct {
		Id struct {
			VideoId string `json:"videoId"`
		} `json:"id"`
		Snippet struct {
			Title        string `json:"title"`
			Description  string `json:"description"`
			ChannelTitle string `json:"channelTitle"`
			PublishedAt  string `json:"publishedAt"`
			Thumbnails   struct {
				Medium struct {
					Url string `json:"url"`
				} `json:"medium"`
			} `json:"thumbnails"`
		} `json:"snippet"`
	} `json:"items"`
}

var channelIDs = []string{
	"UCBwmMxybNva6P_5VmxjzwqA", // UPSC Official
	"UCqYPhGiB9tkShZorfgcL2lA", // SSC Official
	"UC4JX40jDee_tINbkjycV4Sg", // NCS Portal
	"UCL0O2iG0C8XQ6H26o1j7Oaw", // PIB India
	"UC2R2P7Ea4B-7F9aY5FvW_cw", // DD News
	"UCxONzEa1z_6R1E_P0f353qg", // Study IQ Education
	"UCQ-R8O1T7K4O57Q8L12xRnw", // MyGov India
}

var searchQueries = []string{
	"Government Job Notification India",
	"UPSC SSC RRB Latest Vacancy Update",
	"State PSC Job Bharti 2026",
	"Sarkari Online Form Fillup Process",
}

func NewYouTubeSource() *YouTubeSource {
	return &YouTubeSource{
		BaseSource: shared.BaseSource{NameStr: "youtube", BaseURL: "https://www.youtube.com"},
		client: &http.Client{Timeout: 30 * time.Second},
	}
}

func (s *YouTubeSource) Fetch(ctx context.Context) ([]shared.YouTubeVideoSource, error) {
	apiKey := os.Getenv("YOUTUBE_API_KEY")

	if apiKey != "" {
		videos, err := s.fetchAPI(ctx, apiKey)
		if err == nil && len(videos) > 0 {
			return videos, nil
		}
		log.Warn().Err(err).Msg("YouTube API failure/quota reached, falling back to RSS")
	}

	return s.fetchRSS(ctx)
}

func (s *YouTubeSource) fetchAPI(ctx context.Context, key string) ([]shared.YouTubeVideoSource, error) {
	var results []shared.YouTubeVideoSource
	seen := make(map[string]bool)

	for _, q := range searchQueries {
		endpoint := fmt.Sprintf(
			"https://www.googleapis.com/youtube/v3/search?part=snippet&q=%s&type=video&regionCode=IN&relevanceLanguage=hi&maxResults=25&order=date&key=%s",
			url.QueryEscape(q), key)

		resp, err := s.client.Get(endpoint)
		if err != nil {
			continue
		}
		if resp.StatusCode == 403 {
			resp.Body.Close()
			return nil, fmt.Errorf("youtube quota exceeded")
		}

		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		var data YouTubeSearchResponse
		json.Unmarshal(body, &data)

		for _, item := range data.Items {
			if seen[item.Id.VideoId] {
				continue
			}
			seen[item.Id.VideoId] = true

			pubTime, _ := time.Parse(time.RFC3339, item.Snippet.PublishedAt)
			vURL := "https://youtube.com/watch?v=" + item.Id.VideoId

			v := shared.YouTubeVideoSource{
				Source:       "youtube",
				Channel:      item.Snippet.ChannelTitle,
				Title:        item.Snippet.Title,
				URL:          vURL,
				VideoID:      item.Id.VideoId,
				Description:  item.Snippet.Description,
				Thumbnail:    item.Snippet.Thumbnails.Medium.Url,
				Category:     "Education",
				PublishedAt:  &pubTime,
				HashChecksum: s.generateHash(item.Snippet.Title + vURL),
				CreatedAt:    time.Now(),
			}
			results = append(results, v)
		}
	}
	return results, nil
}

func (s *YouTubeSource) fetchRSS(ctx context.Context) ([]shared.YouTubeVideoSource, error) {
	var results []shared.YouTubeVideoSource
	seen := make(map[string]bool)

	for _, cid := range channelIDs {
		feedURL := fmt.Sprintf("https://www.youtube.com/feeds/videos.xml?channel_id=%s", cid)
		resp, err := s.client.Get(feedURL)
		if err != nil || resp.StatusCode != 200 {
			if resp != nil {
				resp.Body.Close()
			}
			continue
		}

		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		doc, err := shared.ParseAtomXML(string(body))
		if err != nil {
			continue
		}

		for _, entry := range doc.Entries {
			if seen[entry.VideoID] {
				continue
			}
			seen[entry.VideoID] = true

			vURL := "https://youtube.com/watch?v=" + entry.VideoID
			pubTime, _ := time.Parse(time.RFC3339, entry.Published)

			v := shared.YouTubeVideoSource{
				Source:       "youtube",
				Channel:      "Official Channel",
				ChannelID:    cid,
				Title:        entry.Title,
				URL:          vURL,
				VideoID:      entry.VideoID,
				Description:  entry.Description,
				Thumbnail:    entry.Thumbnail.URL,
				Category:     "Government",
				PublishedAt:  &pubTime,
				HashChecksum: s.generateHash(entry.Title + vURL),
				CreatedAt:    time.Now(),
			}
			results = append(results, v)
		}
	}
	return results, nil
}

func (s *YouTubeSource) generateHash(input string) string {
	h := sha256.New()
	h.Write([]byte(strings.ToLower(input)))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func (s *YouTubeSource) Name() string {
	return s.NameStr
}
