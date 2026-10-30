UPDATE event_waitlist
SET notified_at = now()
WHERE event_id = $1
  AND notified_at IS NULL
RETURNING
    event_id,
    user_id,
    created_at,
    notified_at;
