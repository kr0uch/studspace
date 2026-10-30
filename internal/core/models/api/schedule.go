package api

import (
	"encoding/json"
	"time"

	"github.com/kr0uch/studspace/internal/core/models/entities"
)

type ScheduleSpace struct {
	SpaceID  int16  `json:"space_id"`
	Name     string `json:"name"`
	Capacity int    `json:"capacity"`
	Color    string `json:"color"`
}

type ScheduleItem struct {
	Type      string                 `json:"type"`
	BookingID int64                  `json:"booking_id,omitempty"`
	EventID   int64                  `json:"event_id,omitempty"`
	SpaceID   int16                  `json:"space_id"`
	Title     string                 `json:"title"`
	StartsAt  time.Time              `json:"starts_at"`
	EndsAt    time.Time              `json:"ends_at"`
	OwnerName string                 `json:"owner_name,omitempty"`
	Organizer string                 `json:"organizer,omitempty"`
	Status    entities.BookingStatus `json:"status,omitempty"`
	IsMine    bool                   `json:"is_mine,omitempty"`
}

type Schedule struct {
	Spaces []ScheduleSpace `json:"spaces"`
	Items  []ScheduleItem  `json:"items"`
}

type BusyDays struct {
	Month time.Time   `json:"month"`
	Days  []time.Time `json:"days"`
}

func (b BusyDays) MarshalJSON() ([]byte, error) {
	days := make([]string, 0, len(b.Days))
	for _, day := range b.Days {
		days = append(days, day.Format("2006-01-02"))
	}
	return json.Marshal(struct {
		Month string   `json:"month"`
		Days  []string `json:"days"`
	}{
		Month: b.Month.Format("2006-01"),
		Days:  days,
	})
}
