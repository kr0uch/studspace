UPDATE users
SET notifications_allowed = $2
WHERE user_id = $1
RETURNING
    user_id,
    vk_user_id,
    first_name,
    last_name,
    photo_url,
    role,
    notifications_allowed,
    created_at;
