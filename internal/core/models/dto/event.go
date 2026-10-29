package dto

import (
	"time"

	"github.com/guregu/null/v6"

	"github.com/kr0uch/studspace/internal/core/models/entities"
)

type CreateEvent struct {
	Title       string
	Organizer   string
	Category    entities.EventCategory
	SpaceID     int16
	StartsAt    time.Time
	EndsAt      time.Time
	Capacity    int
	Description null.String
}

type GetEventViews struct {
	From     time.Time
	To       time.Time
	Category *entities.EventCategory
	Mine     bool
}
