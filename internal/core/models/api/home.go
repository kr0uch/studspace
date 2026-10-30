package api

import (
	"time"

	"github.com/kr0uch/studspace/internal/core/models/entities"
)

type HomeBookingSpace struct {
	SpaceID int16  `json:"space_id"`
	Name    string `json:"name"`
}

type HomeBooking struct {
	BookingID int64                  `json:"booking_id"`
	StartsAt  time.Time              `json:"starts_at"`
	EndsAt    time.Time              `json:"ends_at"`
	Space     HomeBookingSpace       `json:"space"`
	Purpose   string                 `json:"purpose"`
	Status    entities.BookingStatus `json:"status"`
}

type HomeTodayBookings struct {
	Total int           `json:"total"`
	Items []HomeBooking `json:"items"`
}

type Home struct {
	TodayBookings HomeTodayBookings `json:"today_bookings"`
	Spaces        []Space           `json:"spaces"`
}
