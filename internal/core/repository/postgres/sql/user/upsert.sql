INSERT INTO users (vk_user_id, first_name, last_name, photo_url, role)
VALUES ($1, $2, $3, $4, $5::user_role)
ON CONFLICT (vk_user_id) DO UPDATE SET
    first_name = EXCLUDED.first_name,
    last_name = EXCLUDED.last_name,
    photo_url = EXCLUDED.photo_url
RETURNING
    user_id,
    vk_user_id,
    first_name,
    last_name,
    photo_url,
    role,
    notifications_allowed,
    created_at;
