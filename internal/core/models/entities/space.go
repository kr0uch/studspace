package entities

import (
	"time"

	"github.com/guregu/null/v6"
)

type Space struct {
	SpaceID     int16
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
	AmenityID int16
	Name      string
	Icon      string
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
