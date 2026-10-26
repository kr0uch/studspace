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
WHERE space_id = $1
  AND status IN ('pending', 'approved')
  AND period && tstzrange($2, $3, '[)')
  AND ($4::bigint IS NULL OR booking_id <> $4)
ORDER BY starts_at ASC;
