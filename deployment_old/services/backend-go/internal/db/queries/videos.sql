-- name: GetVideos :many
SELECT id, channel, channel_id, title, url, thumbnail, 
       description, video_id, published_at, duration, 
       view_count, like_count, category, created_at
FROM youtube_videos 
WHERE is_active = true
  AND ($1::text = '' OR channel ILIKE '%' || $1 || '%')
  AND ($2::text = '' OR category = $2)
ORDER BY published_at DESC 
LIMIT $3 OFFSET $4;

-- name: GetVideosCount :one
SELECT COUNT(*) FROM youtube_videos WHERE is_active = true
  AND ($1::text = '' OR channel ILIKE '%' || $1 || '%')
  AND ($2::text = '' OR category = $2);

-- name: GetVideoByID :one
SELECT id, channel, channel_id, title, url, thumbnail, 
       description, video_id, published_at, duration, 
       view_count, like_count, category, created_at
FROM youtube_videos 
WHERE id = $1 AND is_active = true;

