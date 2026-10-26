SELECT
    e.event_id,
    e.title,
    e.organizer,
    e.category,
    e.space_id,
    e.starts_at,
    e.ends_at,
    e.capacity,
    e.description,
    e.status,
    e.created_by,
    e.created_at,
    (SELECT count(*) FROM event_registrations r
      WHERE r.event_id = e.event_id AND r.cancelled_at IS NULL) AS registered_count,
    EXISTS (SELECT 1 FROM event_registrations r
      WHERE r.event_id = e.event_id AND r.user_id = $2 AND r.cancelled_at IS NULL) AS is_registered,
    EXISTS (SELECT 1 FROM event_waitlist w
      WHERE w.event_id = e.event_id AND w.user_id = $2 AND w.notified_at IS NULL) AS is_subscribed
FROM events e
WHERE e.event_id = $1;
