package scheduler

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/rojgarsetu/crawler/internal/browser"
	"github.com/rojgarsetu/crawler/internal/proxy"
	"github.com/rojgarsetu/crawler/internal/sources"
	"github.com/rojgarsetu/crawler/internal/store"
)

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

// RunAll instantiates every source, calls Fetch with a per-source timeout,
// routes each result type to the correct store save function, and returns a
// summary. It continues to the next source if one fails.
func RunAll(
	ctx context.Context,
	st *store.PostgresStore,
	browserPool *browser.Pool,
	proxyRotator *proxy.Rotator,
) RunSummary {
	start := time.Now()
	summary := RunSummary{}

	// ---- GovJobFetcher sources (→ jobs_government) ----
	govSources := []sources.GovJobFetcher{
		sources.NewUPSCSource(),
		sources.NewRRBSource(),
		sources.NewSSCSource(),
		sources.NewNCSSource(),
		sources.NewEmploymentNewsSource(),
	}
	for _, s := range govSources {
		summary.SourcesRun++
		summary.SourceResults = append(summary.SourceResults,
			runGovSource(ctx, st, s))
	}

	// ---- PrivJobFetcher sources (→ jobs_private) ----
	privSources := []sources.PrivJobFetcher{
		sources.NewIndeedSource(),
		sources.NewLinkedInSource(),
		sources.NewGoogleJobsSource(),
		sources.NewCompanyPagesSource(),
		sources.NewGreenhouseSource(),
		sources.NewLeverSource(),
	}
	for _, s := range privSources {
		summary.SourcesRun++
		summary.SourceResults = append(summary.SourceResults,
			runPrivSource(ctx, st, s))
	}

	// ---- Naukri (old JobSource interface → convert to PrivJobSource) ----
	summary.SourcesRun++
	summary.SourceResults = append(summary.SourceResults,
		runNaukri(ctx, st, browserPool, proxyRotator))

	// ---- CourseFetcher sources (→ courses) ----
	courseSources := []sources.CourseFetcher{
		sources.NewNPTELSource(),
		sources.NewSWAYAMSource(),
		sources.NewNSDCSource(),
		sources.NewCourseraSource(),
		sources.NewUdemySource(),
	}
	for _, s := range courseSources {
		summary.SourcesRun++
		summary.SourceResults = append(summary.SourceResults,
			runCourseSource(ctx, st, s))
	}

	// ---- VideoFetcher sources (→ youtube_videos) ----
	videoSources := []sources.VideoFetcher{
		sources.NewYouTubeSource(),
	}
	for _, s := range videoSources {
		summary.SourcesRun++
		summary.SourceResults = append(summary.SourceResults,
			runVideoSource(ctx, st, s))
	}

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
	return summary
}

// ── per-type runners ──────────────────────────────────────────────────────────

func runGovSource(ctx context.Context, st *store.PostgresStore, s sources.GovJobFetcher) SourceResult {
	timeout := 3 * time.Minute
	fetchCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	items, err := s.Fetch(fetchCtx)
	if err != nil {
		return SourceResult{Name: s.Name(), Error: fmt.Sprintf("fetch: %v", err)}
	}
	saved := 0
	for i := range items {
		if err := st.SaveGovJob(&items[i]); err != nil {
			log.Printf("ERROR saving gov job from %s: %v", s.Name(), err)
			continue
		}
		saved++
	}
	return SourceResult{Name: s.Name(), Fetched: len(items), Saved: saved}
}

func runPrivSource(ctx context.Context, st *store.PostgresStore, s sources.PrivJobFetcher) SourceResult {
	timeout := 3 * time.Minute
	fetchCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	items, err := s.Fetch(fetchCtx)
	if err != nil {
		return SourceResult{Name: s.Name(), Error: fmt.Sprintf("fetch: %v", err)}
	}
	saved := 0
	for i := range items {
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
	timeout := 5 * time.Minute
	fetchCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	ns := sources.NewNaukriSource(pool, rot)
	items, err := ns.Fetch(fetchCtx)
	if err != nil {
		return SourceResult{Name: name, Error: fmt.Sprintf("fetch: %v", err)}
	}
	saved := 0
	for i := range items {
		// Convert old JobSource → PrivJobSource
		priv := jobSourceToPriv(&items[i])
		if err := st.SavePrivJob(priv); err != nil {
			log.Printf("ERROR saving naukri job: %v", err)
			continue
		}
		saved++
	}
	return SourceResult{Name: name, Fetched: len(items), Saved: saved}
}

// jobSourceToPriv converts the legacy JobSource (used by Naukri) to PrivJobSource.
func jobSourceToPriv(j *sources.JobSource) *sources.PrivJobSource {
	if j == nil {
		return nil
	}
	p := &sources.PrivJobSource{
		Source:      j.Source,
		Company:     j.Company,
		Title:       j.Title,
		Location:    j.Location,
		URL:         j.ApplicationURL,
		Salary:      "",
		Experience:  "",
		JobType:     j.JobType,
		Skills:      j.Skills,
		Description: j.Description,
		PostedAt:    j.PostedAt,
		CreatedAt:   time.Now(),
	}
	if j.SalaryMin != nil && j.SalaryMax != nil {
		p.Salary = fmt.Sprintf("%d-%d", *j.SalaryMin, *j.SalaryMax)
	} else if j.SalaryMin != nil {
		p.Salary = fmt.Sprintf("%d", *j.SalaryMin)
	} else if j.SalaryMax != nil {
		p.Salary = fmt.Sprintf("%d", *j.SalaryMax)
	}
	return p
}

func runCourseSource(ctx context.Context, st *store.PostgresStore, s sources.CourseFetcher) SourceResult {
	timeout := 3 * time.Minute
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

func runVideoSource(ctx context.Context, st *store.PostgresStore, s sources.VideoFetcher) SourceResult {
	timeout := 3 * time.Minute
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
