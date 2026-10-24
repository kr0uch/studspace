package dao

import (
	"fmt"
	"time"

	"github.com/guregu/null/v6"

	"github.com/kr0uch/studspace/internal/core/models/entities"
)

type Space struct {
	ID          int16       `db:"id"`
	Code        string      `db:"code"`
	Name        string      `db:"name"`
	Kind        string      `db:"kind"`
	Capacity    int         `db:"capacity"`
	PhotoURL    null.String `db:"photo_url"`
	Color       string      `db:"color"`
	OpenTime    string      `db:"open_time"`
	CloseTime   string      `db:"close_time"`
	SlotStepMin int         `db:"slot_step_min"`
	IsActive    bool        `db:"is_active"`
}

type Amenity struct {
	ID   int16  `db:"id"`
	Name string `db:"name"`
	Icon string `db:"icon"`
}

type SpaceAmenity struct {
	SpaceID   int16 `db:"space_id"`
	AmenityID int16 `db:"amenity_id"`
}

func (s Space) ToEntity(amenities []Amenity) (entities.Space, error) {
	openTime, err := parseClock(s.OpenTime)
	if err != nil {
		return entities.Space{}, fmt.Errorf("space %d open_time: %w", s.ID, err)
	}
	closeTime, err := parseClock(s.CloseTime)
	if err != nil {
		return entities.Space{}, fmt.Errorf("space %d close_time: %w", s.ID, err)
	}

	result := entities.Space{
		ID:          s.ID,
		Code:        s.Code,
		Name:        s.Name,
		Kind:        s.Kind,
		Capacity:    s.Capacity,
		PhotoURL:    s.PhotoURL,
		Color:       s.Color,
		OpenTime:    openTime,
		CloseTime:   closeTime,
		SlotStepMin: s.SlotStepMin,
		IsActive:    s.IsActive,
		Amenities:   make([]entities.Amenity, 0, len(amenities)),
	}
	for _, a := range amenities {
		result.Amenities = append(result.Amenities, a.ToEntity())
	}
	return result, nil
}

func SpaceFromEntity(e entities.Space) Space {
	return Space{
		ID:          e.ID,
		Code:        e.Code,
		Name:        e.Name,
		Kind:        e.Kind,
		Capacity:    e.Capacity,
		PhotoURL:    e.PhotoURL,
		Color:       e.Color,
		OpenTime:    formatClock(e.OpenTime),
		CloseTime:   formatClock(e.CloseTime),
		SlotStepMin: e.SlotStepMin,
		IsActive:    e.IsActive,
	}
}

func (a Amenity) ToEntity() entities.Amenity {
	return entities.Amenity{
		ID:   a.ID,
		Name: a.Name,
		Icon: a.Icon,
	}
}

func AmenityFromEntity(e entities.Amenity) Amenity {
	return Amenity{
		ID:   e.ID,
		Name: e.Name,
		Icon: e.Icon,
	}
}

func parseClock(value string) (time.Duration, error) {
	t, err := time.Parse(time.TimeOnly, value)
	if err != nil {
		return 0, err
	}
	return time.Duration(t.Hour())*time.Hour +
		time.Duration(t.Minute())*time.Minute +
		time.Duration(t.Second())*time.Second, nil
}

func formatClock(d time.Duration) string {
	total := int(d / time.Second)
	return fmt.Sprintf("%02d:%02d:%02d", total/3600, total%3600/60, total%60)
}
