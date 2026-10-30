SELECT
    event_id,
    title,
    organizer,
    category,
    space_id,
    starts_at,
    ends_at,
    capacity,
    description,
    status,
    created_by,
    created_at
FROM events
WHERE space_id = $1
  AND status = 'published'
  AND period && tstzrange($2, $3, '[)')
  AND ($4::bigint IS NULL OR event_id <> $4)
ORDER BY starts_at ASC, event_id ASC;
