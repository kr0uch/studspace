UPDATE bookings
SET status = 'rejected',
    decided_by = $2,
    decided_at = now(),
    reject_reason_title = $3,
    reject_reason_text = $4
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
