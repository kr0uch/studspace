INSERT INTO booking_status_log (booking_id, from_status, to_status, actor_user_id, comment)
VALUES ($1, $2, $3, $4, $5);
