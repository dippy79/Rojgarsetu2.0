package tests

import (
	"context"
	"database/sql"
	"os"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/rojgarsetu/backend/internal/db"
	"github.com/rojgarsetu/backend/internal/services"
)

func TestSearchAndFiltersPerformance(t *testing.T) {
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
		t.Skipf("Skipping search performance test: db ping failed: %v", err)
		return
	}

	postgresDB := db.NewPostgresDB(sqlDB)
	searchSvc := services.NewSearchService(postgresDB)

	start := time.Now()
	res, err := searchSvc.Search(context.Background(), services.SearchRequest{
		Query: "Engineer",
		Page:  1,
		Limit: 20,
	})
	elapsed := time.Since(start)

	if err != nil {
		t.Fatalf("Search query failed: %v", err)
	}

	if elapsed > 100*time.Millisecond {
		t.Errorf("Search query took too long: %v (expected < 100ms)", elapsed)
	} else {
		t.Logf("Search query executed in %v (GIN index performance PASS). Found %d results.", elapsed, res.Total)
	}
}
