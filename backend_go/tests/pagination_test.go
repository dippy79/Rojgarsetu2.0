package tests

import (
	"database/sql"
	"os"
	"testing"

	_ "github.com/lib/pq"
	"github.com/rojgarsetu/backend/internal/db"
)

func TestPaginationOffsetLogic(t *testing.T) {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@127.0.0.1:5435/rojgarsetu2?sslmode=disable"
	}

	sqlDB, err := sql.Open("postgres", dbURL)
	if err != nil {
		t.Fatalf("Failed to open DB: %v", err)
	}
	defer sqlDB.Close()

	if err := sqlDB.Ping(); err != nil {
		t.Skipf("Skipping pagination test: db ping failed: %v", err)
		return
	}

	postgresDB := db.NewPostgresDB(sqlDB)

	// Page 1
	p1, total1, err1 := postgresDB.GetGovJobs(db.GovJobFilter{}, 1, 5)
	if err1 != nil {
		t.Fatalf("Failed to fetch page 1: %v", err1)
	}

	// Page 2
	p2, total2, err2 := postgresDB.GetGovJobs(db.GovJobFilter{}, 2, 5)
	if err2 != nil {
		t.Fatalf("Failed to fetch page 2: %v", err2)
	}

	if total1 != total2 {
		t.Fatalf("Total count mismatch: %d vs %d", total1, total2)
	}

	if len(p1) != 5 || len(p2) != 5 {
		t.Fatalf("Page size mismatch: len(p1)=%d, len(p2)=%d", len(p1), len(p2))
	}

	// Verify no duplicate IDs between Page 1 and Page 2
	page1IDs := make(map[string]bool)
	for _, job := range p1 {
		page1IDs[job.ID.String()] = true
	}

	for _, job := range p2 {
		if page1IDs[job.ID.String()] {
			t.Fatalf("Pagination Overlap Bug! Job ID %s appeared on both Page 1 and Page 2", job.ID.String())
		}
	}

	t.Logf("Pagination logic PASS. Total: %d, Page 1 and Page 2 returned distinct sets with zero overlap.", total1)
}
