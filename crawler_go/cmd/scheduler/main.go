package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/rojgarsetu/crawler/internal/sources"
	"github.com/rs/zerolog"
)

var (
	logger    zerolog.Logger
	dbURL     string
	redisURL  string
	store     *store.PostgresStore
	scheduler *cron.Cron
)

func main() {
	// Initialize logger
	output := zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}
	logger = zerolog.New(output).With().Timestamp().Logger()

	// Initialize database
	dbURL = os.Getenv("DATABASE_URL")
	if dbURL == "" {
		logger.Fatal().Msg("DATABASE_URL environment variable required")
	}

	var err error
	store, err = store.NewPostgresStore(dbURL)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to connect to database")
	}
	defer store.Close()
	logger.Info().Msg("Database connected")

	// Initialize scheduler
	scheduler = cron.New()

	// Schedule crawlers
	// Government Jobs - every 6 hours
	scheduler.AddFunc("0 */6 * * *", func() {
		logger.Info().Msg("Starting scheduled crawl: Government Jobs")
		runGovJobsCrawler()
	})

	// Private Jobs - every 3 hours
	scheduler.AddFunc("0 */3 * * *", func() {
		logger.Info().Msg("Starting scheduled crawl: Private Jobs")
		runPrivJobsCrawler()
	})

	// Courses - daily
	scheduler.AddFunc("0 0 * * *", func() {
		logger.Info().Msg("Starting scheduled crawl: Courses")
		runCoursesCrawler()
	})

	// YouTube - every 2 hours
	scheduler.AddFunc("0 */2 * * *", func() {
		logger.Info().Msg("Starting scheduled crawl: YouTube")
		runYouTubeCrawler()
	})

	// Start scheduler
	scheduler.Start()
	logger.Info().Msg("Scheduler started")

	// Wait for shutdown signal
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	logger.Info().Msg("Shutting down scheduler...")
	scheduler.Stop()
	logger.Info().Msg("Scheduler stopped")
}

// runGovJobsCrawler crawls government job sources
func runGovJobsCrawler() {
	ctx := context.Background()

	sourcesList := []struct {
		name    string
		fetcher func() sources.GovJobFetcher
	}{
		{"ncs", func() sources.GovJobFetcher { return sources.NewNCSSource() }},
		{"ssc", func() sources.GovJobFetcher { return sources.NewSSCSource() }},
		{"upsc", func() sources.GovJobFetcher { return sources.NewUPSCSource() }},
		{"employment_news", func() sources.GovJobFetcher { return sources.NewEmploymentNewsSource() }},
		{"rrb", func() sources.GovJobFetcher { return sources.NewRRBSource() }},
	}

	totalJobs := 0
	for _, src := range sourcesList {
		fetcher := src.fetcher()
		jobs, err := fetcher.Fetch(ctx)
		if err != nil {
			logger.Error().Err(err).Str("source", src.name).Msg("Failed to fetch government jobs")
			continue
		}

		for _, job := range jobs {
			if err := store.SaveGovJob(job); err != nil {
				logger.Error().Err(err).Str("title", job.Title).Msg("Failed to save job")
			} else {
				totalJobs++
			}
		}

		logger.Info().Int("jobs", len(jobs)).Str("source", src.name).Msg("Government jobs crawled")
	}

	logger.Info().Int("totalJobs", totalJobs).Msg("Government jobs crawl completed")
}

// runPrivJobsCrawler crawls private job sources
func runPrivJobsCrawler() {
	ctx := context.Background()

	sourcesList := []struct {
		name    string
		fetcher func() sources.PrivJobFetcher
	}{
		{"linkedin", func() sources.PrivJobFetcher { return sources.NewLinkedInSource() }},
		{"indeed", func() sources.PrivJobFetcher { return sources.NewIndeedSource() }},
		{"google_jobs", func() sources.PrivJobFetcher { return sources.NewGoogleJobsSource() }},
		{"company_pages", func() sources.PrivJobFetcher { return sources.NewCompanyPagesSource() }},
	}

	totalJobs := 0
	for _, src := range sourcesList {
		fetcher := src.fetcher()
		jobs, err := fetcher.Fetch(ctx)
		if err != nil {
			logger.Error().Err(err).Str("source", src.name).Msg("Failed to fetch private jobs")
			continue
		}

		for _, job := range jobs {
			if err := store.SavePrivJob(job); err != nil {
				logger.Error().Err(err).Str("title", job.Title).Msg("Failed to save job")
			} else {
				totalJobs++
			}
		}

		logger.Info().Int("jobs", len(jobs)).Str("source", src.name).Msg("Private jobs crawled")
	}

	logger.Info().Int("totalJobs", totalJobs).Msg("Private jobs crawl completed")
}

// runCoursesCrawler crawls course sources
func runCoursesCrawler() {
	ctx := context.Background()

	sourcesList := []struct {
		name    string
		fetcher func() sources.CourseFetcher
	}{
		{"nptel", func() sources.CourseFetcher { return sources.NewNPTELSource() }},
		{"swayam", func() sources.CourseFetcher { return sources.NewSWAYAMSource() }},
		{"nsdc", func() sources.CourseFetcher { return sources.NewNSDCSource() }},
		{"coursera", func() sources.CourseFetcher { return sources.NewCourseraSource() }},
		{"udemy", func() sources.CourseFetcher { return sources.NewUdemySource() }},
	}

	totalCourses := 0
	for _, src := range sourcesList {
		fetcher := src.fetcher()
		courses, err := fetcher.Fetch(ctx)
		if err != nil {
			logger.Error().Err(err).Str("source", src.name).Msg("Failed to fetch courses")
			continue
		}

		for _, course := range courses {
			if err := store.SaveCourse(course); err != nil {
				logger.Error().Err(err).Str("title", course.Title).Msg("Failed to save course")
			} else {
				totalCourses++
			}
		}

		logger.Info().Int("courses", len(courses)).Str("source", src.name).Msg("Courses crawled")
	}

	logger.Info().Int("totalCourses", totalCourses).Msg("Courses crawl completed")
}

// runYouTubeCrawler crawls YouTube videos
func runYouTubeCrawler() {
	ctx := context.Background()

	fetcher := sources.NewYouTubeSource()
	videos, err := fetcher.Fetch(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to fetch YouTube videos")
		return
	}

	totalVideos := 0
	for _, video := range videos {
		if err := store.SaveYouTubeVideo(video); err != nil {
			logger.Error().Err(err).Str("title", video.Title).Msg("Failed to save video")
		} else {
			totalVideos++
		}
	}

	logger.Info().Int("totalVideos", totalVideos).Msg("YouTube crawl completed")
}

// Initialize log
func init() {
	log.SetOutput(os.Stdout)
}
