\echo === COLUMNS CHECK ===
SELECT table_name, column_name
FROM information_schema.columns
WHERE table_schema='public'
  AND column_name='language'
ORDER BY table_name;

\echo === LANGUAGE DISTRIBUTION ===
SELECT 'jobs_government' AS tbl, language, COUNT(*) 
FROM jobs_government GROUP BY language
UNION ALL
SELECT 'jobs_private', language, COUNT(*) 
FROM jobs_private GROUP BY language
UNION ALL
SELECT 'company_jobs', language, COUNT(*) 
FROM company_jobs GROUP BY language
UNION ALL
SELECT 'courses', language, COUNT(*) 
FROM courses GROUP BY language
UNION ALL
SELECT 'youtube_videos', language, COUNT(*) 
FROM youtube_videos GROUP BY language
ORDER BY tbl, language;
