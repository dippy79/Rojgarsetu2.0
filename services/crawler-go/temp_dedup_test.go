package main

import (
	"os"
	"testing"
	"time"

	"github.com/rojgarsetu/crawler/internal/shared"
	"github.com/rojgarsetu/crawler/internal/store"
)

func TestCrawlerDeduplication(t *testing.T) {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@localhost:5435/rojgarsetu2?sslmode=disable"
	}

	st, err := store.NewPostgresStore(dbURL)
	if err != nil {
		t.Skipf("Skipping crawler dedup test: database connection failed: %v", err)
		return
	}
	defer st.Close()

	// 1. Sample Data
	govJob := &shared.GovJobSource{
		Title:      "Dedup Test Assistant Officer",
		Department: "Staff Selection Commission",
		Location:   "New Delhi",
		ApplyURL:   "https://ssc.gov.in/job/dedup123",
		Source:     "SSC",
	}

	privJob := &shared.PrivJobSource{
		Company:     "Dedup Tech Solutions",
		Title:       "Dedup Lead Software Engineer",
		Location:    "Bangalore",
		URL:         "https://deduptech.example.com/job/456",
		Source:      "Direct",
		Skills:      []string{"Go", "PostgreSQL"},
		Description: "High performance backend role",
	}

	course := &shared.CourseSource{
		Provider: "NPTEL",
		Title:    "Advanced Go Concurrency Dedup",
		URL:      "https://nptel.ac.in/courses/dedup_go_1",
		Source:   "NPTEL",
		Skills:   []string{"Go", "Concurrency"},
	}

	now := time.Now()
	video := &shared.YouTubeVideoSource{
		Channel:     "Tech Education",
		ChannelID:   "UC_dedup_123",
		Title:       "Mastering Go Pipelines",
		URL:         "https://youtube.com/watch?v=dedup_vid_999",
		Thumbnail:   "https://img.youtube.com/vi/dedup_vid_999/0.jpg",
		Description: "Deep dive into Go channels and mutexes",
		VideoID:     "dedup_vid_999",
		PublishedAt: &now,
		Category:    "Engineering",
	}

	// 2. First Run: Save All Items
	if err := st.SaveGovJob(govJob); err != nil {
		t.Fatalf("Failed 1st SaveGovJob: %v", err)
	}
	if err := st.SavePrivJob(privJob); err != nil {
		t.Fatalf("Failed 1st SavePrivJob: %v", err)
	}
	if err := st.SaveCourse(course); err != nil {
		t.Fatalf("Failed 1st SaveCourse: %v", err)
	}
	if err := st.SaveVideo(video); err != nil {
		t.Fatalf("Failed 1st SaveVideo: %v", err)
	}

	var countGov1, countPriv1, countCourse1, countVideo1 int
	_ = st.DB().QueryRow("SELECT COUNT(*) FROM jobs_government WHERE title = $1", govJob.Title).Scan(&countGov1)
	_ = st.DB().QueryRow("SELECT COUNT(*) FROM jobs_private WHERE title = $1", privJob.Title).Scan(&countPriv1)
	_ = st.DB().QueryRow("SELECT COUNT(*) FROM courses WHERE url = $1", course.URL).Scan(&countCourse1)
	_ = st.DB().QueryRow("SELECT COUNT(*) FROM youtube_videos WHERE video_id = $1", video.VideoID).Scan(&countVideo1)

	if countGov1 != 1 || countPriv1 != 1 || countCourse1 != 1 || countVideo1 != 1 {
		t.Fatalf("First run expected 1 record per table, got gov=%d, priv=%d, course=%d, video=%d",
			countGov1, countPriv1, countCourse1, countVideo1)
	}

	// 3. Second Run: Save Exact Same Items Again
	if err := st.SaveGovJob(govJob); err != nil {
		t.Fatalf("Failed 2nd SaveGovJob: %v", err)
	}
	if err := st.SavePrivJob(privJob); err != nil {
		t.Fatalf("Failed 2nd SavePrivJob: %v", err)
	}
	if err := st.SaveCourse(course); err != nil {
		t.Fatalf("Failed 2nd SaveCourse: %v", err)
	}
	if err := st.SaveVideo(video); err != nil {
		t.Fatalf("Failed 2nd SaveVideo: %v", err)
	}

	var countGov2, countPriv2, countCourse2, countVideo2 int
	_ = st.DB().QueryRow("SELECT COUNT(*) FROM jobs_government WHERE title = $1", govJob.Title).Scan(&countGov2)
	_ = st.DB().QueryRow("SELECT COUNT(*) FROM jobs_private WHERE title = $1", privJob.Title).Scan(&countPriv2)
	_ = st.DB().QueryRow("SELECT COUNT(*) FROM courses WHERE url = $1", course.URL).Scan(&countCourse2)
	_ = st.DB().QueryRow("SELECT COUNT(*) FROM youtube_videos WHERE video_id = $1", video.VideoID).Scan(&countVideo2)

	if countGov2 != 1 || countPriv2 != 1 || countCourse2 != 1 || countVideo2 != 1 {
		t.Fatalf("Deduplication failed! Second run produced duplicate rows: gov=%d, priv=%d, course=%d, video=%d",
			countGov2, countPriv2, countCourse2, countVideo2)
	}

	t.Log("Crawler deduplication verified 100% PASS on all tables!")
}
