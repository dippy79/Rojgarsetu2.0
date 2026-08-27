package sources

import (
	"context"
	"github.com/rojgarsetu/crawler/internal/videos"
)

type YouTubeScraper struct {
	BaseSource
	internalSource *videos.YouTubeSource
}

func NewYouTubeScraper() *YouTubeScraper {
	return &YouTubeScraper{
		BaseSource:     BaseSource{NameStr: "youtube"},
		internalSource: videos.NewYouTubeSource(),
	}
}

func (s *YouTubeScraper) FetchVideos(ctx context.Context) ([]YouTubeVideoSource, error) {
	internalVideos, err := s.internalSource.Fetch(ctx)
	if err != nil {
		return nil, err
	}

	// Direct map as they are same underlying struct
	return internalVideos, nil
}

func (s *YouTubeScraper) Name() string {
	return s.NameStr
}
