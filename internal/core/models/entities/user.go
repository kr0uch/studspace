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
	ID                   int64
	VKUserID             int64
	FirstName            string
	LastName             string
	PhotoURL             null.String
	Role                 UserRole
	NotificationsAllowed bool
	CreatedAt            time.Time
	LastSeenAt           null.Time
}

func (u User) IsAdmin() bool {
	return u.Role == UserRoleAdmin
}

func (u User) FullName() string {
	return u.FirstName + " " + u.LastName
}
