INSERT INTO events (
    title,
    organizer,
    category,
    space_id,
    starts_at,
    ends_at,
    capacity,
    description,
    status,
    created_by
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
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
