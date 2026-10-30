UPDATE bookings
SET status = 'approved',
    decided_by = $2,
    decided_at = now(),
    reject_reason_title = NULL,
    reject_reason_text = NULL,
    cancelled_by = NULL,
    cancelled_at = NULL
WHERE booking_id = $1
RETURNING
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
    updated_at;
