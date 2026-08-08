# Videos Page Bug Fix — TODO

## Root causes (evidence gathered)
1. **"All Channels" dropdown empty** — `VideosPage.jsx` hardcodes `channels: []`; no `/api/v1/videos/channels` endpoint exists (404).
2. **Category dropdown "does nothing"** — options hardcoded (`Jobs/Education/Skills/Tech/Interviews`) don't match real DB categories (`Government/Tech Skills/Career Prep/Data Science`), so selected values return empty results.
3. **One-per-page instead of grid** — client-side `isGovVideo` filter discards the full newest-page of Government videos, leaving 0-1 tech videos. Backend `limit=20` is correct.

## Backend steps
- [x] 1. `backend_go/internal/db/database.go`: add `GetVideoChannels()` + `GetVideoCategories()` (distinct raw queries)
- [x] 2. `backend_go/internal/db/database.go`: add exclusion support to `GetVideos` (excludeCategory) so Tech-tab count is post-exclusion
- [x] 3. `backend_go/internal/services/video_service.go`: expose new methods + exclusion passthrough
- [x] 4. `backend_go/internal/handlers/video_handler.go`: add `GetVideoChannels`/`GetVideoCategories` handlers; accept `exclude` query in `GetVideos`
- [x] 5. `backend_go/cmd/server/main.go`: register `/api/v1/videos/channels` + `/api/v1/videos/categories` routes

## Frontend steps
- [x] 6. `frontend/src/pages/videos/VideosPage.jsx`: fetch channels + categories from new endpoints; populate FilterBar with real data
- [x] 7. `VideosPage.jsx`: Tech tab sends `exclude=Government` to API; remove client-side `isGovVideo` filter

## Verification
- [x] 8. Raw response from `/api/v1/videos/channels` (8 channels)
- [x] 9. Raw response from `/api/v1/videos/categories` (4 categories)
- [x] 10. Tech tab call with `exclude=Government` — confirm list AND pagination total are post-exclusion
- [x] 11. Browser dropdown/grid reflects the fix
