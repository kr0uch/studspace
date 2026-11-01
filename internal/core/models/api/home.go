package api

import (
	"time"

	"github.com/kr0uch/studspace/internal/core/models/entities"
)

const (
	HomeItemBooking = "booking"
	HomeItemEvent   = "event"
)

type HomeTodaySpace struct {
	SpaceID int16  `json:"space_id"`
	Name    string `json:"name"`
}

type HomeTodayItem struct {
	Type      string                 `json:"type"`
	BookingID int64                  `json:"booking_id,omitempty"`
	EventID   int64                  `json:"event_id,omitempty"`
	StartsAt  time.Time              `json:"starts_at"`
	EndsAt    time.Time              `json:"ends_at"`
	Space     HomeTodaySpace         `json:"space"`
	Purpose   string                 `json:"purpose,omitempty"`
	Title     string                 `json:"title,omitempty"`
	Status    entities.BookingStatus `json:"status,omitempty"`
}

type HomeToday struct {
	Total int             `json:"total"`
	Items []HomeTodayItem `json:"items"`
}

type Home struct {
	Today  HomeToday `json:"today"`
	Spaces []Space   `json:"spaces"`
}
