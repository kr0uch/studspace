UPDATE events
SET title = $2,
    organizer = $3,
    category = $4,
    space_id = $5,
    starts_at = $6,
    ends_at = $7,
    capacity = $8,
    description = $9
WHERE event_id = $1
  AND status = 'published'
RETURNING
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
    created_at;
