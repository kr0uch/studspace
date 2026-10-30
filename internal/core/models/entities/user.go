package entities

import (
	"time"

	"github.com/guregu/null/v6"
)

type UserRole string

const (
	UserRoleStudent UserRole = "student"
	UserRoleAdmin   UserRole = "admin"
)

type User struct {
	UserID               int64       `json:"user_id"`
	VKUserID             int64       `json:"vk_user_id"`
	FirstName            string      `json:"first_name"`
	LastName             string      `json:"last_name"`
	PhotoURL             null.String `json:"photo_url"`
	Role                 UserRole    `json:"role"`
	NotificationsAllowed bool        `json:"notifications_allowed"`
	CreatedAt            time.Time   `json:"-"`
}

func (u User) IsAdmin() bool {
	return u.Role == UserRoleAdmin
}

func (u User) FullName() string {
	return u.FirstName + " " + u.LastName
}
