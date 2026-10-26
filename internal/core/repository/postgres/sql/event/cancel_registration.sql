UPDATE event_registrations
SET cancelled_at = now()
WHERE event_id = $1
  AND user_id = $2
  AND cancelled_at IS NULL;
