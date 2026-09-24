package tests

import (
	"database/sql"
	"os"
	"sync"
	"testing"

	_ "github.com/lib/pq"
	"github.com/rojgarsetu/backend/config"
)

func TestDBConnectionPooling(t *testing.T) {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@localhost:5435/rojgarsetu2?sslmode=disable"
	}

	cfg := config.Load()
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		t.Fatalf("Failed to open db: %v", err)
	}
	defer db.Close()

	db.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	db.SetMaxIdleConns(cfg.Database.MaxIdleConns)

	if err := db.Ping(); err != nil {
		t.Skipf("Skipping DB pooling test: postgres ping failed: %v", err)
		return
	}

	var wg sync.WaitGroup
	concurrentCount := 50
	errChan := make(chan error, concurrentCount)

	for i := 0; i < concurrentCount; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			var result int
			err := db.QueryRow("SELECT 1").Scan(&result)
			if err != nil {
				errChan <- err
			}
		}(i)
	}

	wg.Wait()
	close(errChan)

	var errorsEncountered []error
	for err := range errChan {
		errorsEncountered = append(errorsEncountered, err)
	}

	if len(errorsEncountered) > 0 {
		t.Fatalf("Encountered %d errors during concurrent pooling test. First error: %v", len(errorsEncountered), errorsEncountered[0])
	}

	t.Logf("Successfully executed %d concurrent DB queries with zero pooling errors.", concurrentCount)
}
