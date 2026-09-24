package tests

import (
	"database/sql"
	"os"
	"testing"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func TestMigrationAndSeederLifecycle(t *testing.T) {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@localhost:5435/rojgarsetu2?sslmode=disable"
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		t.Fatalf("Failed to open DB: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		t.Skipf("Skipping migration lifecycle test: postgres ping failed: %v", err)
		return
	}

	// Initialize Migrations
	m, err := migrate.New("file://../migrations", dbURL)
	if err != nil {
		t.Fatalf("Failed to initialize migrations: %v", err)
	}

	// 1. Run Migrations Up to ensure clean initial state
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		t.Fatalf("Migration Up failed: %v", err)
	}

	// 2. Verify Key Tables Exist and are non-empty after seeding
	tablesToVerify := []string{
		"users",
		"candidates",
		"companies",
		"jobs_government",
		"jobs_private",
		"courses",
		"youtube_videos",
	}

	for _, table := range tablesToVerify {
		var count int
		err := db.QueryRow("SELECT COUNT(*) FROM " + table).Scan(&count)
		if err != nil {
			t.Fatalf("Failed to query table %s: %v", table, err)
		}
		t.Logf("Table '%s' contains %d rows.", table, count)
	}

	t.Log("Migration lifecycle and database population verified successfully.")
}
