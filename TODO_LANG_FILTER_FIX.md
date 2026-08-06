# Language Filter Bug Fix

## Root Cause
`backend_go/internal/handlers/search_handler.go` `Search()` (POST handler):
- `ShouldBindJSON(&req)` correctly sets `req.Language` from JSON body
- Then `req.Language = c.Query("language")` overwrites it with empty query-string value
- Empty string → SQL `($1 = '' OR language = $1)` → always TRUE → matches everything

## Plan Steps
- [x] Fix 1: Remove `req.Language = c.Query("language")` from Search() POST handler (body first, query fallback)
- [x] Fix 2: Harden AI-engine `preferred_language` filter in recommender/service.py (SQL pushdown + defensive check)
- [x] Rebuild backend and AI engine
- [x] Verify: search language:"en" vs "hi" produce different counts
- [x] Verify: recommendations preferred_language:"en" vs "hi" produce different counts

## Verification Results (live, after --force-recreate + --build --no-cache)
- Backend `/api/v1/search` language:"en" => total 206; language:"hi" => total 0  (BEFORE fix: both 206)
- AI `/recommend/jobs` preferred_language:"en" => total 10; preferred_language:"hi" => 0 (no Hindi jobs)
- DB: priv_en=817, priv_hi=0, search_en=206, search_hi=0
