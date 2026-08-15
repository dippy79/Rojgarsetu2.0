package crawler

import (
    "crypto/sha256"
    "encoding/hex"
    "fmt"
    "strings"
)

// GenerateJobHash creates a unique SHA-256 fingerprint for deduplication
func GenerateJobHash(sourceName, title, organization, applyURL string) string {
    cleanTitle := strings.ToLower(strings.TrimSpace(title))
    cleanOrg := strings.ToLower(strings.TrimSpace(organization))
    cleanURL := strings.ToLower(strings.TrimSpace(applyURL))
    cleanSource := strings.ToLower(strings.TrimSpace(sourceName))

    raw := fmt.Sprintf("%s|%s|%s|%s", cleanSource, cleanTitle, cleanOrg, cleanURL)
    hash := sha256.Sum256([]byte(raw))
    return hex.EncodeToString(hash[:])
}
