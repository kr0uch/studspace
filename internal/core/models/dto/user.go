package dto

import "github.com/guregu/null/v6"

type UpsertUserByVK struct {
	VKUserID  int64
	FirstName string
	LastName  string
	PhotoURL  null.String
}

type UpdateUserNotifications struct {
	UserID  int64
	Allowed bool
}
