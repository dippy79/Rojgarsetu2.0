package scheduler

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/rojgarsetu/crawler/internal/browser"
	"github.com/rojgarsetu/crawler/internal/courses"
	"github.com/rojgarsetu/crawler/internal/jobs/gov"
	"github.com/rojgarsetu/crawler/internal/jobs/priv"
	"github.com/rojgarsetu/crawler/internal/legal"
	"github.com/rojgarsetu/crawler/internal/proxy"
	"github.com/rojgarsetu/crawler/internal/shared"
	"github.com/rojgarsetu/crawler/internal/sources"
	"github.com/rojgarsetu/crawler/internal/store"
	"github.com/rojgarsetu/crawler/internal/videos"
)

var antiFake = legal.NewAntiFakeEngine()

// RunSummary captures the aggregate result of a crawl run.
type RunSummary struct {
	SourcesRun    int            `json:"sources_run"`
	Succeeded     int            `json:"succeeded"`
	Failed        int            `json:"failed"`
	TotalSaved    int            `json:"total_saved"`
	SourceResults []SourceResult `json:"source_results"`
	Duration      string         `json:"duration"`
}

// SourceResult captures per-source outcome.
type SourceResult struct {
	Name    string `json:"name"`
	Fetched int    `json:"fetched"`
	Saved   int    `json:"saved"`
	Error   string `json:"error,omitempty"`
}

func RunAll(
	ctx context.Context,
	st *store.PostgresStore,
	browserPool *browser.Pool,
	proxyRotator *proxy.Rotator,
) RunSummary {
	start := time.Now()
	summary := RunSummary{}

	// ---- GovJobFetcher sources (→ jobs_government) ----
	govSources := []shared.GovJobFetcher{
		gov.NewUPSCSource(browserPool),
		gov.NewRRBSource(browserPool),
		gov.NewSSCSource(browserPool),
		gov.NewNCSSource(browserPool),
	}
	for _, s := range govSources {
		summary.SourcesRun++
		summary.SourceResults = append(summary.SourceResults, runGovSource(ctx, st, s))
	}

	// ---- State Jobs (New) ----
	stateJobs := sources.NewStateJobsSource(browserPool)
	summary.SourcesRun++
	summary.SourceResults = append(summary.SourceResults, runGovSource(ctx, st, stateJobs))

	// ---- PrivJobFetcher sources (→ jobs_private) ----
	privSources := []shared.PrivJobFetcher{
		priv.NewIndeedSource(),
		priv.NewLinkedInSource(),
		priv.NewGoogleJobsSource(),
		priv.NewCompanyPagesSource(),
		priv.NewGreenhouseSource(),
		priv.NewLeverSource(),
		priv.NewApnaSource(),
		priv.NewInternshalaSource(),
		priv.NewShineSource(),
	}
	for _, s := range privSources {
		summary.SourcesRun++
		summary.SourceResults = append(summary.SourceResults, runPrivSource(ctx, st, s))
	}

	// ---- Naukri (old JobSource interface → convert to PrivJobSource) ----
	summary.SourcesRun++
	summary.SourceResults = append(summary.SourceResults, runNaukri(ctx, st, browserPool, proxyRotator))

	// ---- CourseFetcher sources (→ courses) ----
	courseSources := []shared.CourseFetcher{
		courses.NewNPTELSource(),
		courses.NewSWAYAMSource(),
		courses.NewNSDCSource(),
		courses.NewCourseraSource(),
		courses.NewUdemySource(),
		courses.NewGeeksforGeeksSource(),
		courses.NewTutorialsPointSource(),
		courses.NewW3SchoolsSource(),
	}
	for _, s := range courseSources {
		summary.SourcesRun++
		summary.SourceResults = append(summary.SourceResults, runCourseSource(ctx, st, s))
	}

	// ---- VideoFetcher sources (→ youtube_videos) ----
	videoSources := []shared.VideoFetcher{
		videos.NewYouTubeSource(),
	}
	for _, s := range videoSources {
		summary.SourcesRun++
		summary.SourceResults = append(summary.SourceResults, runVideoSource(ctx, st, s))
	}

	// ---- Forms Scraper (New) ----
	formsScraper := sources.NewGovFormsSource(browserPool)
	summary.SourcesRun++
	summary.SourceResults = append(summary.SourceResults, runFormsSource(ctx, st, formsScraper))

	// ---- Aggregate totals ----
	for _, r := range summary.SourceResults {
		summary.TotalSaved += r.Saved
		if r.Error != "" {
			summary.Failed++
		} else {
			summary.Succeeded++
		}
	}
	summary.Duration = time.Since(start).Round(time.Second).String()

	// Log to database
	status := "COMPLETED"
	errMsg := ""
	if summary.Failed > 0 {
		status = "PARTIAL_FAILURE"
		errMsg = fmt.Sprintf("%d sources failed", summary.Failed)
	}
	if err := st.SaveLog(summary.SourcesRun, 0, summary.TotalSaved, 0, status, errMsg); err != nil {
		log.Printf("ERROR saving crawl log: %v", err)
	}

	return summary
}

func runGovSource(ctx context.Context, st *store.PostgresStore, s shared.GovJobFetcher) SourceResult {
	timeout := 5 * time.Minute
	fetchCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	items, err := s.Fetch(fetchCtx)
	if err != nil {
		return SourceResult{Name: s.Name(), Error: fmt.Sprintf("fetch: %v", err)}
	}
	saved := 0
	for i := range items {
		// Anti-Fake Validation
		verified, reason := antiFake.ValidateGovJob(items[i].Title, items[i].Department, items[i].ApplyURL)
		items[i].IsVerified = verified
		items[i].VerificationMeta = map[string]any{
			"engine": "AntiFakeEngine v2",
			"reason": reason,
			"ts":     time.Now().Unix(),
		}
		if !verified {
			items[i].ScamScore = 0.8
		}

		if err := st.SaveGovJob(&items[i]); err != nil {
			log.Printf("ERROR saving gov job from %s: %v", s.Name(), err)
			continue
		}
		saved++
	}
	return SourceResult{Name: s.Name(), Fetched: len(items), Saved: saved}
}

func runPrivSource(ctx context.Context, st *store.PostgresStore, s shared.PrivJobFetcher) SourceResult {
	timeout := 5 * time.Minute
	fetchCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	items, err := s.Fetch(fetchCtx)
	if err != nil {
		return SourceResult{Name: s.Name(), Error: fmt.Sprintf("fetch: %v", err)}
	}
	saved := 0
	for i := range items {
		// Anti-Fake Validation (Basic for Private)
		isScam := antiFake.ContainsScamKeywords(items[i].Title) || antiFake.ContainsScamKeywords(items[i].Description)
		items[i].IsVerified = !isScam
		if isScam {
			items[i].ScamScore = 0.9
		}
		items[i].VerificationMeta = map[string]any{
			"engine": "AntiFakeEngine v2",
			"isScam": isScam,
			"ts":     time.Now().Unix(),
		}

		if err := st.SavePrivJob(&items[i]); err != nil {
			log.Printf("ERROR saving priv job from %s: %v", s.Name(), err)
			continue
		}
		saved++
	}
	return SourceResult{Name: s.Name(), Fetched: len(items), Saved: saved}
}

func runNaukri(ctx context.Context, st *store.PostgresStore, pool *browser.Pool, rot *proxy.Rotator) SourceResult {
	name := "naukri"
	timeout := 8 * time.Minute
	fetchCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	ns := priv.NewNaukriSource(pool, rot)
	items, err := ns.Fetch(fetchCtx)
	if err != nil {
		return SourceResult{Name: name, Error: fmt.Sprintf("fetch: %v", err)}
	}
	saved := 0
	for i := range items {
		privItem := jobSourceToPriv(&items[i])
		if err := st.SavePrivJob(privItem); err != nil {
			log.Printf("ERROR saving naukri job: %v", err)
			continue
		}
		saved++
	}
	return SourceResult{Name: name, Fetched: len(items), Saved: saved}
}

func jobSourceToPriv(j *shared.JobSource) *shared.PrivJobSource {
	if j == nil {
		return nil
	}
	p := &shared.PrivJobSource{
		Source:      j.Source,
		Company:     j.Company,
		Title:       j.Title,
		Location:    j.Location,
		URL:         j.ApplicationURL,
		JobType:     j.JobType,
		Skills:      j.Skills,
		Description: j.Description,
		PostedAt:    j.PostedAt,
		CreatedAt:   time.Now(),
	}
	return p
}

func runCourseSource(ctx context.Context, st *store.PostgresStore, s shared.CourseFetcher) SourceResult {
	timeout := 5 * time.Minute
	fetchCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	items, err := s.Fetch(fetchCtx)
	if err != nil {
		return SourceResult{Name: s.Name(), Error: fmt.Sprintf("fetch: %v", err)}
	}
	saved := 0
	for i := range items {
		if err := st.SaveCourse(&items[i]); err != nil {
			log.Printf("ERROR saving course from %s: %v", s.Name(), err)
			continue
		}
		saved++
	}
	return SourceResult{Name: s.Name(), Fetched: len(items), Saved: saved}
}

func runVideoSource(ctx context.Context, st *store.PostgresStore, s shared.VideoFetcher) SourceResult {
	timeout := 5 * time.Minute
	fetchCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	items, err := s.Fetch(fetchCtx)
	if err != nil {
		return SourceResult{Name: s.Name(), Error: fmt.Sprintf("fetch: %v", err)}
	}
	saved := 0
	for i := range items {
		if err := st.SaveVideo(&items[i]); err != nil {
			log.Printf("ERROR saving video from %s: %v", s.Name(), err)
			continue
		}
		saved++
	}
	return SourceResult{Name: s.Name(), Fetched: len(items), Saved: saved}
}

func runFormsSource(ctx context.Context, st *store.PostgresStore, s *sources.GovFormsSource) SourceResult {
	timeout := 8 * time.Minute
	fetchCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	items, err := s.Fetch(fetchCtx)
	if err != nil {
		return SourceResult{Name: s.Name(), Error: fmt.Sprintf("fetch: %v", err)}
	}
	saved := 0
	for i := range items {
		if err := st.SaveGovForm(&items[i]); err != nil {
			log.Printf("ERROR saving gov form from %s: %v", s.Name(), err)
			continue
		}
		saved++
	}
	return SourceResult{Name: s.Name(), Fetched: len(items), Saved: saved}
}
