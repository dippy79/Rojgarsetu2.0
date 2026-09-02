package priv

import (
	"os"
	"testing"
)

// TestParseHTMLJobs_JSONLDFallback verifies that when regex patterns find no
// jobs, the JSON-LD fallback extracts schema.org JobPosting entries from the
// embedded application/ld+json blocks.
func TestParseHTMLJobs_JSONLDFallback(t *testing.T) {
	s := NewCompanyPagesSource()

	html := `<html><head>
<script type="application/ld+json">
[{"@context":"https://schema.org","@type":"JobPosting","title":"Software Engineer","description":"<p>Build stuff</p>","datePosted":"2024-01-02T00:00:00Z","hiringOrganization":{"name":"Acme Corp"},"jobLocation":{"@type":"Place","address":{"addressLocality":"Bangalore","addressRegion":"KA","addressCountry":"IN"}},"baseSalary":{"currency":"INR","value":{"minValue":1000000,"maxValue":2000000}}}]
</script>
</head><body></body></html>`

	jobs := s.parseHTMLJobs(html, "Acme", "https://acme.example/careers")
	if len(jobs) != 1 {
		t.Fatalf("expected 1 job from JSON-LD, got %d", len(jobs))
	}
	job := jobs[0]
	if job.Title != "Software Engineer" {
		t.Errorf("title mismatch: %q", job.Title)
	}
	if job.Company != "Acme Corp" {
		t.Errorf("company mismatch: %q", job.Company)
	}
	if job.Location != "Bangalore, KA" {
		t.Errorf("location mismatch: %q", job.Location)
	}
	if job.Salary == "" {
		t.Errorf("expected salary to be populated")
	}
	if job.Description == "" {
		t.Errorf("expected description to be populated")
	}
}

// TestParseHTMLJobs_NoJSONLD_ExistingCompanies verifies the honest no-op
// expectation: the raw HTML captured for the existing 15-company list exposes
// no JSON-LD JobPosting data, so the parser returns 0 for them.
func TestParseHTMLJobs_NoJSONLD_ExistingCompanies(t *testing.T) {
	s := NewCompanyPagesSource()

	cases := []struct {
		name string
		path string
	}{
		{"tcs", "tcs.html"},
		{"wipro", "wipro.html"},
	}

	for _, c := range cases {
		data, err := os.ReadFile(c.path)
		if err != nil {
			t.Skipf("fixture %s not present: %v", c.path, err)
			continue
		}
		jobs := s.parseHTMLJobs(string(data), c.name, "https://"+c.name+".example/careers")
		if len(jobs) != 0 {
			t.Errorf("%s: expected 0 jobs from raw HTML (JSON-LD absent), got %d", c.name, len(jobs))
		}
	}
}
