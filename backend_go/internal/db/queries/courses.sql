-- name: GetCourses :many
SELECT id, provider, title, url, duration, mode, level, skills,
       description, thumbnail_url, price, is_free, source, language, created_at
FROM courses 
WHERE is_active = true
  AND ($1::text = '' OR provider ILIKE '%' || $1 || '%')
  AND ($2::text = '' OR mode = $2)
  AND ($3::text = '' OR level = $3)
  AND ($4::text = '' OR language = $4)
ORDER BY created_at DESC 
LIMIT $5 OFFSET $6;

-- name: GetCoursesCount :one
SELECT COUNT(*) FROM courses WHERE is_active = true
  AND ($1::text = '' OR provider ILIKE '%' || $1 || '%')
  AND ($2::text = '' OR mode = $2)
  AND ($3::text = '' OR level = $3)
  AND ($4::text = '' OR language = $4);

-- name: GetCourseByID :one
SELECT id, provider, title, url, duration, mode, level, skills,
       description, thumbnail_url, price, is_free, source, language, created_at
FROM courses 
WHERE id = $1 AND is_active = true;

