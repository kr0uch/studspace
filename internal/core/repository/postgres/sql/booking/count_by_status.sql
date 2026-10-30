SELECT count(*)
FROM bookings
WHERE status = $1::booking_status;
