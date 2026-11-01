UPDATE users
SET patronymic = $2,
    occupation = $3
WHERE user_id = $1
RETURNING
    user_id,
    vk_user_id,
    first_name,
    last_name,
    patronymic,
    occupation,
    photo_url,
    role,
    notifications_allowed,
    created_at;
