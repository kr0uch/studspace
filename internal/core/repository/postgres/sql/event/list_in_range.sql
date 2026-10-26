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
WHERE status = 'published'
  AND period && tstzrange($1, $2, '[)')
  AND ($3::smallint IS NULL OR space_id = $3)
ORDER BY starts_at ASC, event_id ASC;
