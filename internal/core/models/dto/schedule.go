package dto

import (
	"time"

	"github.com/kr0uch/studspace/internal/core/models/entities"
)

type Home struct {
	TodayBookings []entities.Booking
	Spaces        []SpaceWithNow
}

type SpaceWithNow struct {
	Space entities.Space
	Now   *entities.SpaceNowStatus
}

type ScheduleItemKind string

const (
	ScheduleItemBooking ScheduleItemKind = "booking"
	ScheduleItemEvent   ScheduleItemKind = "event"
)

type ScheduleItem struct {
	Kind    ScheduleItemKind
	Booking *entities.Booking
	Event   *entities.Event
	IsMine  bool
}

type GetSchedule struct {
	From time.Time
	To   time.Time
}

type Schedule struct {
	Spaces []entities.Space
	Items  []ScheduleItem
}
