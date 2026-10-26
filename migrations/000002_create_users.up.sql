CREATE TABLE users (
    user_id               BIGSERIAL PRIMARY KEY,
    vk_user_id            BIGINT      NOT NULL UNIQUE,
    first_name            TEXT        NOT NULL,
    last_name             TEXT        NOT NULL,
    photo_url             TEXT,
    role                  user_role   NOT NULL DEFAULT 'student',
    notifications_allowed BOOLEAN     NOT NULL DEFAULT false,
    created_at            TIMESTAMPTZ NOT NULL DEFAULT now()
);
