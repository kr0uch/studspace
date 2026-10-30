SELECT
    booking_id,
    user_id,
    space_id,
    first_name,
    last_name,
    starts_at,
    ends_at,
    purpose,
    people_count,
    status,
    reject_reason_title,
    reject_reason_text,
    decided_by,
    decided_at,
    cancelled_by,
    cancelled_at,
    rescheduled_from_id,
    created_at,
    updated_at
FROM bookings
WHERE user_id = $1
  AND NOT (status = 'pending' OR (status = 'approved' AND ends_at > $2))
  AND ($3::timestamptz IS NULL OR (starts_at, booking_id) < ($3, $4))
ORDER BY starts_at DESC, booking_id DESC
LIMIT $5;
