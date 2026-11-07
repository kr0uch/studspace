SELECT EXISTS (
    SELECT 1
    FROM bookings
    WHERE rescheduled_from_id = $1
      AND status = 'pending'
);
