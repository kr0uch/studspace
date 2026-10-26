SELECT count(*)
FROM bookings
WHERE user_id = $1
  AND status = 'pending';
