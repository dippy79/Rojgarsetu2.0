package crawler

import (
    "crypto/sha256"
    "fmt"
    "strings"
)

// GenerateJobHash creates a SHA-256 hash for job deduplication
func GenerateJobHash(title, applyURL string) string {
    input := fmt.Sprintf("%s|%s",
        strings.ToLower(strings.TrimSpace(title)),
        strings.TrimSpace(applyURL))
    hash := sha256.Sum256([]byte(input))
    return fmt.Sprintf("%x", hash)
}