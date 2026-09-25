package main

import (
	"database/sql"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	_ "github.com/lib/pq"
)

func TestExternalMediaFeedsAndCourses(t *testing.T) {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@127.0.0.1:5435/rojgarsetu2?sslmode=disable"
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		t.Fatalf("Failed to open DB: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		t.Skipf("Skipping media feeds test: db ping failed: %v", err)
		return
	}

	client := &http.Client{
		Timeout: 5 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= 5 {
				return http.ErrUseLastResponse
			}
			return nil
		},
	}

	// 1. Verify Video Links
	rowsVideos, err := db.Query("SELECT id, url FROM youtube_videos WHERE is_active = true LIMIT 10")
	if err != nil {
		t.Fatalf("Failed to query videos: %v", err)
	}
	defer rowsVideos.Close()

	var validVideos, totalVideos int
	for rowsVideos.Next() {
		var id, videoURL string
		if err := rowsVideos.Scan(&id, &videoURL); err != nil {
			continue
		}
		totalVideos++

		// Ensure URL has protocol
		if !strings.HasPrefix(videoURL, "http://") && !strings.HasPrefix(videoURL, "https://") {
			videoURL = "https://" + videoURL
			_, _ = db.Exec("UPDATE youtube_videos SET url = $1 WHERE id = $2", videoURL, id)
		}

		req, err := http.NewRequest("HEAD", videoURL, nil)
		if err != nil {
			continue
		}
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64)")
		resp, err := client.Do(req)
		if err == nil {
			resp.Body.Close()
			if resp.StatusCode >= 200 && resp.StatusCode < 400 {
				validVideos++
			}
		}
	}

	// 2. Verify Course Links
	rowsCourses, err := db.Query("SELECT id, url FROM courses WHERE is_active = true LIMIT 10")
	if err != nil {
		t.Fatalf("Failed to query courses: %v", err)
	}
	defer rowsCourses.Close()

	var validCourses, totalCourses int
	for rowsCourses.Next() {
		var id, courseURL string
		if err := rowsCourses.Scan(&id, &courseURL); err != nil {
			continue
		}
		totalCourses++

		if !strings.HasPrefix(courseURL, "http://") && !strings.HasPrefix(courseURL, "https://") {
			courseURL = "https://" + courseURL
			_, _ = db.Exec("UPDATE courses SET url = $1 WHERE id = $2", courseURL, id)
		}

		req, err := http.NewRequest("HEAD", courseURL, nil)
		if err != nil {
			continue
		}
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64)")
		resp, err := client.Do(req)
		if err == nil {
			resp.Body.Close()
			if resp.StatusCode >= 200 && resp.StatusCode < 400 {
				validCourses++
			}
		}
	}

	t.Logf("Media Feed Verification: %d/%d Videos active, %d/%d Courses active.", validVideos, totalVideos, validCourses, totalCourses)
}
