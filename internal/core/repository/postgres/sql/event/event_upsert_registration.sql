INSERT INTO event_registrations (event_id, user_id)
VALUES ($1, $2)
ON CONFLICT (event_id, user_id) DO UPDATE
SET cancelled_at = NULL,
    created_at = now()
WHERE event_registrations.cancelled_at IS NOT NULL;
