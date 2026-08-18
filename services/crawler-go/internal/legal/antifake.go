package legal

import (
	"crypto/sha256"
	"fmt"
	"net/url"
	"strings"
)

// AntiFakeEngine provides validation for job listings
type AntiFakeEngine struct {
	scamKeywords []string
	govDomains   []string
}

// NewAntiFakeEngine creates a new AntiFakeEngine
func NewAntiFakeEngine() *AntiFakeEngine {
	return &AntiFakeEngine{
		scamKeywords: []string{
			"payment required", "deposit money", "processing fee",
			"whatsapp for job", "pay to join", "security deposit",
			"registration fee", "fake job", "scam",
		},
		govDomains: []string{
			".gov.in", ".nic.in", ".ac.in", ".edu.in", ".res.in",
		},
	}
}

// IsTrustedGovDomain checks if the URL belongs to a trusted government domain
func (e *AntiFakeEngine) IsTrustedGovDomain(applyURL string) bool {
	u, err := url.Parse(applyURL)
	if err != nil {
		return false
	}
	host := strings.ToLower(u.Host)
	for _, domain := range e.govDomains {
		if strings.HasSuffix(host, domain) {
			return true
		}
	}
	return false
}

// ContainsScamKeywords checks if the text contains any known scam keywords
func (e *AntiFakeEngine) ContainsScamKeywords(text string) bool {
	lowerText := strings.ToLower(text)
	for _, kw := range e.scamKeywords {
		if strings.Contains(lowerText, kw) {
			return true
		}
	}
	return false
}

// GenerateSHA256Hash creates a SHA256 hash for deduplication
func GenerateSHA256Hash(title, company, url string) string {
	data := fmt.Sprintf("%s|%s|%s", strings.ToLower(title), strings.ToLower(company), strings.ToLower(url))
	hash := sha256.Sum256([]byte(data))
	return fmt.Sprintf("%x", hash)
}

// ValidateGovJob performs full validation for a government job
func (e *AntiFakeEngine) ValidateGovJob(title, dept, applyURL string) (bool, string) {
	if !e.IsTrustedGovDomain(applyURL) {
		return false, "Unverified domain for government job"
	}
	if e.ContainsScamKeywords(title) || e.ContainsScamKeywords(dept) {
		return false, "Potential scam keywords detected"
	}
	return true, "Verified"
}
