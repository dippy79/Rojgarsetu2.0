# PART A: React Web Frontend — Implementation TODO

Target: `frontend/` (React). Backend UNTOUCHED. ✅ COMPLETED

- [x] Phase 1: Courses — dynamic providers, all-provider sequential list
- [x] Phase 2: Videos — Tech/News tab split, default Tech-only view
- [x] Phase 3: Navbar — Jobs dropdown (Gov/Private), Courses, Videos + expose job_region & language filters
- [x] Phase 4: CSS revamp — per-screen gradients + glassmorphism
- [x] Verify: `npm run build` compiles without errors

## Files Modified
- `src/pages/courses/CoursesPage.jsx` — dynamic providers, graceful `(courses || [])` loading
- `src/pages/videos/VideosPage.jsx` — Tech default + Government/News secondary tab
- `src/components/Navbar.jsx` — 🏢 Jobs dropdown (Gov/Private), 📚 Courses, 🎥 Videos
- `src/components/FilterBar.jsx` — added `regions` + `languages` filter groups
- `src/pages/gov-jobs/GovJobsPage.jsx` — region/language wired
- `src/pages/private-jobs/PrivateJobsPage.jsx` — region/language wired
- CSS revamp (all 11 files): index.css, Navbar.css, JobCard.css, FilterBar.css, Pagination.css, CourseCard.css, VideoCard.css, GovJobs.css, PrivateJobs.css, Courses.css, Videos.css

## Verification
- `npm run build` → **"Compiled successfully." / "The build folder is ready to be deployed."**
- Backend + Flutter app untouched.

