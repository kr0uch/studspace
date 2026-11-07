package dto

import (
	"time"

	"github.com/kr0uch/studspace/internal/core/models/entities"
)

type BusyInterval struct {
	StartsAt time.Time
	EndsAt   time.Time
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
