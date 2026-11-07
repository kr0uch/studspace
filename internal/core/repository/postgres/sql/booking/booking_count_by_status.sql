SELECT count(*)
FROM bookings
WHERE status = $1::booking_status
  AND (status <> 'pending' OR starts_at > $2);
