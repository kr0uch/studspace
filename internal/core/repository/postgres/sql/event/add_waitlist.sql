INSERT INTO event_waitlist (event_id, user_id)
VALUES ($1, $2)
ON CONFLICT (event_id, user_id) DO NOTHING;
