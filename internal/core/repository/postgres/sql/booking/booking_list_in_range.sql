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
WHERE period && tstzrange($1, $2, '[)')
  AND ($3::smallint IS NULL OR space_id = $3)
  AND status IN ('pending', 'approved')
ORDER BY starts_at ASC, booking_id ASC;
