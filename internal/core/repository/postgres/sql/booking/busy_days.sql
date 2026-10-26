SELECT DISTINCT timezone('Europe/Moscow', starts_at)::date::text AS day
FROM bookings
WHERE status IN ('pending', 'approved')
  AND starts_at >= $1
  AND starts_at < $2
ORDER BY day;
