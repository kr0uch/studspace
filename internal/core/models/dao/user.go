package dao

import (
	"time"

	"github.com/guregu/null/v6"

	"github.com/kr0uch/studspace/internal/core/models/entities"
)

type User struct {
	ID                   int64       `db:"id"`
	VKUserID             int64       `db:"vk_user_id"`
	FirstName            string      `db:"first_name"`
	LastName             string      `db:"last_name"`
	PhotoURL             null.String `db:"photo_url"`
	Role                 string      `db:"role"`
	NotificationsAllowed bool        `db:"notifications_allowed"`
	CreatedAt            time.Time   `db:"created_at"`
	LastSeenAt           null.Time   `db:"last_seen_at"`
}

func (u User) ToEntity() entities.User {
	return entities.User{
		ID:                   u.ID,
		VKUserID:             u.VKUserID,
		FirstName:            u.FirstName,
		LastName:             u.LastName,
		PhotoURL:             u.PhotoURL,
		Role:                 entities.UserRole(u.Role),
		NotificationsAllowed: u.NotificationsAllowed,
		CreatedAt:            u.CreatedAt,
		LastSeenAt:           u.LastSeenAt,
	}
}

func UserFromEntity(e entities.User) User {
	return User{
		ID:                   e.ID,
		VKUserID:             e.VKUserID,
		FirstName:            e.FirstName,
		LastName:             e.LastName,
		PhotoURL:             e.PhotoURL,
		Role:                 string(e.Role),
		NotificationsAllowed: e.NotificationsAllowed,
		CreatedAt:            e.CreatedAt,
		LastSeenAt:           e.LastSeenAt,
	}
}
