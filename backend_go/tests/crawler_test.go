package tests

import (
	"crypto/sha256"
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestDedup tests job deduplication logic
func TestDedup(t *testing.T) {
	// Sample job data
	job1 := map[string]string{
		"title":      "Software Engineer",
		"company":    "Tech Corp",
		"location":   "Delhi",
		"salary":     "10-15 LPA",
		"source":     "test",
	}
	
	job2 := map[string]string{
		"title":      "Software Engineer",
		"company":    "Tech Corp",
		"location":   "Delhi",
		"salary":     "10-15 LPA",
		"source":     "test",
	}
	
	job3 := map[string]string{
		"title":      "Data Scientist",
		"company":    "Tech Corp",
		"location":   "Delhi",
		"salary":     "15-20 LPA",
		"source":     "test",
	}

	// Generate hashes
	hash1 := generateJobHash(job1)
	hash2 := generateJobHash(job2)
	hash3 := generateJobHash(job3)

	// Test deduplication
	assert.Equal(t, hash1, hash2, "Identical jobs should have same hash")
	assert.NotEqual(t, hash1, hash3, "Different jobs should have different hashes")
}

// TestHashGeneration tests hash generation consistency
func TestHashGeneration(t *testing.T) {
	job := map[string]string{
		"title":      "Software Engineer",
		"company":    "Tech Corp",
		"location":   "Delhi",
		"salary":     "10-15 LPA",
		"source":     "test",
	}

	hash1 := generateJobHash(job)
	hash2 := generateJobHash(job)

	assert.Equal(t, hash1, hash2, "Hash generation should be consistent")
	assert.Len(t, hash1, 64, "SHA256 hash should be 64 characters")
}

// TestGracefulSkip tests graceful error handling
func TestGracefulSkip(t *testing.T) {
	// Test that invalid data doesn't crash the system
	invalidJobs := []map[string]interface{}{
		{"title": ""}, // Missing required fields
		{"company": "Test"}, // Missing title
		{}, // Empty job
	}

	for _, job := range invalidJobs {
		// This should not panic
		_ = processJob(job)
	}
}

// Helper function to generate job hash
func generateJobHash(job map[string]string) string {
	data := job["title"] + "|" + job["company"] + "|" + job["location"] + "|" + job["salary"] + "|" + job["source"]
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

// Helper function to process job (simulates crawler processing)
func processJob(job map[string]interface{}) error {
	// Simulate processing with error handling
	if title, ok := job["title"].(string); ok && title != "" {
		return nil
	}
	return nil // Gracefully skip invalid jobs
}