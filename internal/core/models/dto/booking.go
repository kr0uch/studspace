package dto

import (
	"time"

	"github.com/kr0uch/studspace/internal/core/models/entities"
)

type CreateBooking struct {
	SpaceID     int16
	StartsAt    time.Time
	EndsAt      time.Time
	Purpose     string
	PeopleCount int
}

type RescheduleBooking struct {
	BookingID int64
	SpaceID   int16
	StartsAt  time.Time
	EndsAt    time.Time
}

const (
	DefaultListLimit = 20
	MaxListLimit     = 100
)

type ListUserBookings struct {
	UserID int64
	Tab    entities.UserBookingTab
	Now    time.Time
	Cursor string
	Limit  int
}

type ListAdminBookings struct {
	Tab    entities.AdminBookingTab
	Cursor string
	Limit  int
}

type GetSlotOverlaps struct {
	SpaceID          int16
	StartsAt         time.Time
	EndsAt           time.Time
	ExcludeBookingID *int64
}

type SlotOverlaps struct {
	Pending  []entities.Booking
	Approved []entities.Booking
	Events   []entities.Event
}

func (o SlotOverlaps) HardConflict() bool {
	return len(o.Approved) > 0 || len(o.Events) > 0
}
