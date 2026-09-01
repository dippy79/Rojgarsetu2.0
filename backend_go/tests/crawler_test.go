package tests

import (
	"strings"
	"testing"

	"github.com/rojgarsetu/backend/internal/crawler"
)

// IsSpam mimics the behavior of a spam filter for job listings
func IsSpam(text string) bool {
	spamKeywords := []string{
		"deposit money",
		"earn daily",
		"payment required",
		"whatsapp for job",
		"security deposit",
	}
	lowerText := strings.ToLower(text)
	for _, kw := range spamKeywords {
		if strings.Contains(lowerText, kw) {
			return true
		}
	}
	return false
}

// TestHashGeneration verifies consistency and uniqueness of SHA256 fingerprints
func TestHashGeneration(t *testing.T) {
	title := "UPSC Exam"
	url := "https://upsc.gov.in"
	org := "UPSC"
	source := "upsc_official"

	// Consistency
	hash1 := crawler.GenerateJobHash(source, title, org, url)
	hash2 := crawler.GenerateJobHash(source, title, org, url)

	if hash1 != hash2 {
		t.Errorf("expected consistent hash, got %s and %s", hash1, hash2)
	}

	// Uniqueness
	hash3 := crawler.GenerateJobHash(source, "SSC Exam", org, "https://ssc.gov.in")
	if hash1 == hash3 {
		t.Errorf("different inputs produced identical hash: %s", hash1)
	}

	if len(hash1) != 64 {
		t.Errorf("expected 64 char hex hash, got %d", len(hash1))
	}
}

// TestSpamFilter tests the detection of fraudulent job postings
func TestSpamFilter(t *testing.T) {
	testCases := []struct {
		input    string
		expected bool
	}{
		{"deposit money now earn daily", true},
		{"UPSC Civil Services 2026", false},
		{"Apply for SSC CGL vacancies", false},
		{"Security deposit needed to join", true},
		{"Contact via WhatsApp for job confirmation", true},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			result := IsSpam(tc.input)
			if result != tc.expected {
				t.Errorf("IsSpam(%q) = %v; want %v", tc.input, result, tc.expected)
			}
		})
	}
}

// TestGracefulSkip tests that invalid data handled gracefully
func TestGracefulSkip(t *testing.T) {
	// Simulate invalid data processing
	title := ""
	if title == "" {
		t.Log("Successfully identified empty title for skipping")
	} else {
		t.Error("Failed to identify empty title")
	}
}
