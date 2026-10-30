SELECT DISTINCT timezone('Europe/Moscow', starts_at)::date AS day
FROM events
WHERE status = 'published'
  AND starts_at >= $1
  AND starts_at < $2
ORDER BY day;
