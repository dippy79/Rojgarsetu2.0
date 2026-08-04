// Package lang provides a lightweight, dependency-free language detection
// heuristic for the RojgarSetu crawler. It is intentionally NOT a full NLP
// library (no external deps, no vendored model weights). It scores the
// presence of Indic script Unicode ranges plus a small set of English-latin
// signals and returns an ISO 639-1 code.
//
// Design goals:
//   - Robust for the common RojgarSetu case where title/description are
//     English-heavy and often have NO strong script signal (e.g. Naukri,
//     Greenhouse, Lever, Coursera). In that case we MUST confidently return
//     "en", NOT misdetect, and NOT return "auto" (which is a detection mode,
//     not a stored value).
//   - Pure Unicode codepoint scoring, so it works on any UTF-8 input without
//     tokenizers or dictionaries.
//   - conservative: only returns a non-English code when a script's character
//     count exceeds a clear threshold (so a single Devanagari character in an
//     otherwise-English description cannot flip the result).
package lang

import (
	"strings"
	"unicode"
)

// Detect returns an ISO 639-1 language code for the given text.
//
// It examines the concatenation of title + description (and any optional
// extra fields). If the text contains no strong script signal (e.g. empty
// description, pure ASCII, or numbers/URLs only), it returns "en".
//
// The threshold logic: we count the number of runes in each recognized script
// block. A non-English language is returned only if its script count is BOTH
// >= 3 AND >= 20% of the total non-space runes. This guarantees the common
// English-heavy row with a stray Indic character still resolves to "en".
func Detect(title, description string, extra ...string) string {
	parts := make([]string, 0, 2+len(extra))
	parts = append(parts, title, description)
	parts = append(parts, extra...)
	text := strings.Join(parts, " ")

	// Count non-space runes so we can compute a ratio for the confidence gate.
	var totalRunes int
	scriptCounts := map[string]int{}
	for _, r := range text {
		if unicode.IsSpace(r) {
			continue
		}
		totalRunes++
		if code := scriptForRune(r); code != "" {
			scriptCounts[code]++
		}
	}

	// No meaningful content at all (empty title+description, or only
	// whitespace/punctuation) → confidently English. This is the common
	// Naukri/Greenhouse/Lever row with an empty description.
	if totalRunes == 0 {
		return "en"
	}

	// Strong script signal threshold: require at least 3 characters AND at
	// least 20% of the total non-space runes. This prevents a single stray
	// Indic character from flipping an English-heavy row.
	best := ""
	bestCount := 0
	for code, count := range scriptCounts {
		if count >= 3 && count*5 >= totalRunes { // count >= totalRunes/5
			if count > bestCount {
				best = code
				bestCount = count
			}
		}
	}
	if best != "" {
		return best
	}

	// No strong script signal → default to English (matches the column
	// default and the real, overwhelmingly-English corpus).
	return "en"
}

// scriptForRune maps a rune to an ISO 639-1 language code if it falls in a
// recognized Indic script block. Latin letters are deliberately NOT mapped to
// any language here — they are the implicit default (English) when no Indic
// script dominates.
func scriptForRune(r rune) string {
	switch {
	// Devanagari covers Hindi, Marathi, Nepali, Sanskrit.
	case isDevanagari(r):
		return "hi"
	case isTamil(r):
		return "ta"
	case isTelugu(r):
		return "te"
	// Bengali covers Bengali and Assamese.
	case isBengali(r):
		return "bn"
	case isGujarati(r):
		return "gu"
	case isKannada(r):
		return "kn"
	case isMalayalam(r):
		return "ml"
	case isPunjabi(r):
		return "pa"
	case isUrdu(r):
		return "ur"
	case isOriya(r):
		return "or"
	}
	return ""
}

// ---- Unicode block helpers ------------------------------------------------

// Devanagari: U+0900–U+097F (Hindi, Marathi, Nepali, Sanskrit)
func isDevanagari(r rune) bool {
	return r >= 0x0900 && r <= 0x097F
}

// Tamil: U+0B80–U+0BFF
func isTamil(r rune) bool {
	return r >= 0x0B80 && r <= 0x0BFF
}

// Telugu: U+0C00–U+0C7F
func isTelugu(r rune) bool {
	return r >= 0x0C00 && r <= 0x0C7F
}

// Bengali: U+0980–U+09FF (Bengali, Assamese)
func isBengali(r rune) bool {
	return r >= 0x0980 && r <= 0x09FF
}

// Gujarati: U+0A80–U+0AFF
func isGujarati(r rune) bool {
	return r >= 0x0A80 && r <= 0x0AFF
}

// Kannada: U+0C80–U+0CFF
func isKannada(r rune) bool {
	return r >= 0x0C80 && r <= 0x0CFF
}

// Malayalam: U+0D00–U+0D7F
func isMalayalam(r rune) bool {
	return r >= 0x0D00 && r <= 0x0D7F
}

// Punjabi (Gurmukhi): U+0A00–U+0A7F
func isPunjabi(r rune) bool {
	return r >= 0x0A00 && r <= 0x0A7F
}

// Urdu (Arabic script): U+0600–U+06FF
func isUrdu(r rune) bool {
	return r >= 0x0600 && r <= 0x06FF
}

// Oriya: U+0B00–U+0B7F
func isOriya(r rune) bool {
	return r >= 0x0B00 && r <= 0x0B7F
}
