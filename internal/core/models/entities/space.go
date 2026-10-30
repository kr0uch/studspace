package entities

import (
	"time"

	"github.com/guregu/null/v6"
)

type Space struct {
	SpaceID     int16         `json:"space_id"`
	Code        string        `json:"code"`
	Name        string        `json:"name"`
	Kind        string        `json:"kind"`
	Capacity    int           `json:"capacity"`
	PhotoURL    null.String   `json:"photo_url"`
	Color       string        `json:"color"`
	OpenTime    time.Duration `json:"-"`
	CloseTime   time.Duration `json:"-"`
	SlotStepMin int           `json:"-"`
	IsActive    bool          `json:"-"`
	Amenities   []Amenity     `json:"amenities"`
}

type Amenity struct {
	AmenityID int16  `json:"amenity_id"`
	Name      string `json:"name"`
	Icon      string `json:"icon"`
}

type SpaceState string

const (
	SpaceStateFree   SpaceState = "free"
	SpaceStateBusy   SpaceState = "busy"
	SpaceStateClosed SpaceState = "closed"
)

type SpaceNowStatus struct {
	State SpaceState `json:"state"`
	Until null.Time  `json:"until"`
}
