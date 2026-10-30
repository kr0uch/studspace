UPDATE events
SET status = $2::event_status
WHERE event_id = $1
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
