package dto

import (
	"time"

	"github.com/kr0uch/studspace/internal/core/models/entities"
)

type BusyKind string

const (
	BusyKindBooking BusyKind = "booking"
	BusyKindEvent   BusyKind = "event"
)

type BusyInterval struct {
	StartsAt time.Time
	EndsAt   time.Time
	Kind     BusyKind
}

type GetSpaceAvailability struct {
	SpaceID int16
	Date    time.Time
}

type GetSpaceCurrentStatus struct {
	SpaceID int16
	Now     time.Time
}

type SpaceAvailability struct {
	Space entities.Space
	Date  time.Time
	Busy  []BusyInterval
}
