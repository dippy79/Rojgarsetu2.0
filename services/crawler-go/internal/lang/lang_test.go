package lang

import "testing"

// TestDetectEmptyDescription verifies the critical default behavior: when a
// row has a title but an EMPTY description (the common Naukri/Greenhouse/
// Lever/Coursera case where the crawler could not extract a description), the
// detector must return "en" — NOT a misdetection, NOT "auto".
func TestDetectEmptyDescription(t *testing.T) {
	cases := []struct {
		name        string
		title       string
		description string
		want        string
	}{
		{
			name:        "empty title and empty description",
			title:       "",
			description: "",
			want:        "en",
		},
		{
			name:        "english title, empty description",
			title:       "Software Engineer",
			description: "",
			want:        "en",
		},
		{
			name:        "english title, whitespace-only description",
			title:       "Backend Developer",
			description: "   \n  ",
			want:        "en",
		},
		{
			name:        "english title, url-only description (no script signal)",
			title:       "Frontend Engineer",
			description: "https://boards.greenhouse.io/stripe/jobs/123",
			want:        "en",
		},
		{
			name:        "english title, numbers-only description",
			title:       "Data Analyst",
			description: "12345 67890",
			want:        "en",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := Detect(tc.title, tc.description)
			if got != tc.want {
				t.Fatalf("Detect(%q, %q) = %q, want %q", tc.title, tc.description, got, tc.want)
			}
		})
	}
}

// TestDetectStrongIndicSignal verifies that a genuinely non-English row with a
// strong Indic script signal IS detected, so we don't over-tune the detector
// into always returning "en".
func TestDetectStrongIndicSignal(t *testing.T) {
	cases := []struct {
		name        string
		title       string
		description string
		want        string
	}{
		{
			name:        "hindi title + hindi description",
			title:       "सॉफ्टवेयर इंजीनियर",
			description: "नौकरी के लिए आवेदन करें",
			want:        "hi",
		},
		{
			name:        "tamil description",
			title:       "வேலை",
			description: "இது ஒரு தமிழ் வேலை விளக்கம்",
			want:        "ta",
		},
		{
			name:        "english title but strong hindi description",
			title:       "Software Engineer",
			description: "यह एक पूर्णकालिक नौकरी है जिसमें आपको आवेदन करना होगा",
			want:        "hi",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := Detect(tc.title, tc.description)
			if got != tc.want {
				t.Fatalf("Detect(%q, %q) = %q, want %q", tc.title, tc.description, got, tc.want)
			}
		})
	}
}

// TestDetectStrayIndicCharEnglishHeavy verifies the conservative threshold:
// a single stray Devanagari character in an otherwise-English description must
// NOT flip the result to "hi".
func TestDetectStrayIndicCharEnglishHeavy(t *testing.T) {
	title := "Software Engineer at Stripe"
	desc := "This is a full-time software engineering role at Stripe. " +
		"Responsibilities include designing, building, and maintaining " +
		"distributed systems. Requires strong experience in Go and Python. " +
		"Check out the requirements section for more details. हा"
	got := Detect(title, desc)
	if got != "en" {
		t.Fatalf("Detect() = %q, want %q (stray Indic char must not flip result)", got, "en")
	}
}
