package api

import "github.com/guregu/null/v6"

type SyncUserRequest struct {
	FirstName string      `json:"first_name"`
	LastName  string      `json:"last_name"`
	PhotoURL  null.String `json:"photo_url" swaggertype:"string"`
}

type PatchNotificationsRequest struct {
	Allowed bool `json:"allowed"`
}

type PatchMeRequest struct {
	Patronymic null.String `json:"patronymic" swaggertype:"string"`
	Occupation null.String `json:"occupation" swaggertype:"string"`
}
