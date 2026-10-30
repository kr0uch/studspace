package api

import "github.com/guregu/null/v6"

type SyncUserRequest struct {
	FirstName string      `json:"first_name"`
	LastName  string      `json:"last_name"`
	PhotoURL  null.String `json:"photo_url"`
}

type PatchNotificationsRequest struct {
	Allowed bool `json:"allowed"`
}
