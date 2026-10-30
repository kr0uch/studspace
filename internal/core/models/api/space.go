package api

import (
	"time"

	"github.com/kr0uch/studspace/internal/core/models/dto"
	"github.com/kr0uch/studspace/internal/core/models/entities"
	"github.com/kr0uch/studspace/pkg/utils"
)

type Space struct {
	entities.Space
	Now *entities.SpaceNowStatus `json:"now,omitempty"`
}

type BusyInterval struct {
	StartsAt time.Time `json:"starts_at"`
	EndsAt   time.Time `json:"ends_at"`
	Type     string    `json:"type"`
}

type SpaceAvailability struct {
	SpaceID     int16          `json:"space_id"`
	Date        string         `json:"date"`
	OpenTime    string         `json:"open_time"`
	CloseTime   string         `json:"close_time"`
	SlotStepMin int            `json:"slot_step_min"`
	Busy        []BusyInterval `json:"busy"`
}

func NewSpaceAvailability(in dto.SpaceAvailability) SpaceAvailability {
	busy := make([]BusyInterval, 0, len(in.Busy))
	for _, interval := range in.Busy {
		busy = append(busy, BusyInterval{
			StartsAt: interval.StartsAt,
			EndsAt:   interval.EndsAt,
			Type:     string(interval.Kind),
		})
	}
	return SpaceAvailability{
		SpaceID:     in.Space.SpaceID,
		Date:        in.Date.In(utils.Moscow).Format("2006-01-02"),
		OpenTime:    utils.FormatClock(in.Space.OpenTime),
		CloseTime:   utils.FormatClock(in.Space.CloseTime),
		SlotStepMin: in.Space.SlotStepMin,
		Busy:        busy,
	}
}
