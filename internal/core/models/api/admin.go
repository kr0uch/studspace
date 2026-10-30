package api

import (
	"time"

	"github.com/kr0uch/studspace/internal/core/models/entities"
)

type AdminBookingConflict struct {
	BookingID int64                  `json:"booking_id"`
	OwnerName string                 `json:"owner_name"`
	StartsAt  time.Time              `json:"starts_at"`
	EndsAt    time.Time              `json:"ends_at"`
	Status    entities.BookingStatus `json:"status"`
}

type AdminBooking struct {
	Booking
	User      *entities.User         `json:"user"`
	Conflicts []AdminBookingConflict `json:"conflicts"`
}

type AdminBookingCounters struct {
	New int `json:"new"`
}

type AdminBookingList struct {
	Counters   AdminBookingCounters `json:"counters"`
	Items      []AdminBooking       `json:"items"`
	NextCursor *string              `json:"next_cursor"`
}
