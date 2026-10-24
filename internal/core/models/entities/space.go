package entities

import (
	"time"

	"github.com/guregu/null/v6"
)

type Space struct {
	ID          int16
	Code        string
	Name        string
	Kind        string
	Capacity    int
	PhotoURL    null.String
	Color       string
	OpenTime    time.Duration
	CloseTime   time.Duration
	SlotStepMin int
	IsActive    bool
	Amenities   []Amenity
}

type Amenity struct {
	ID   int16
	Name string
	Icon string
}

type SpaceState string

const (
	SpaceStateFree SpaceState = "free"
	SpaceStateBusy SpaceState = "busy"
)

type SpaceNowStatus struct {
	State SpaceState
	Until null.Time
}

type OccupancyType string

const (
	OccupancyTypeBooking OccupancyType = "booking"
	OccupancyTypeEvent   OccupancyType = "event"
)

type Occupancy struct {
	SpaceID  int16
	Type     OccupancyType
	SourceID int64
	StartsAt time.Time
	EndsAt   time.Time
}

func (o Occupancy) Overlaps(startsAt, endsAt time.Time) bool {
	return o.StartsAt.Before(endsAt) && startsAt.Before(o.EndsAt)
}
