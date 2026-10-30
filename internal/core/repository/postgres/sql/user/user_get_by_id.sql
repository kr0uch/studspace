SELECT
    user_id,
    vk_user_id,
    first_name,
    last_name,
    photo_url,
    role,
    notifications_allowed,
    created_at
FROM users
WHERE user_id = $1;
