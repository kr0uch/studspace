package api

import (
	"time"

	"github.com/guregu/null/v6"

	"github.com/kr0uch/studspace/internal/core/models/entities"
)

type CreateEventRequest struct {
	Title       string                 `json:"title"`
	Organizer   string                 `json:"organizer"`
	Category    entities.EventCategory `json:"category"`
	SpaceID     int16                  `json:"space_id"`
	StartsAt    time.Time              `json:"starts_at"`
	EndsAt      time.Time              `json:"ends_at"`
	Capacity    int                    `json:"capacity"`
	Description null.String            `json:"description"`
}

type EventSpace struct {
	SpaceID int16  `json:"space_id"`
	Name    string `json:"name,omitempty"`
}

type Event struct {
	EventID         int64                   `json:"event_id"`
	Title           string                  `json:"title"`
	Organizer       string                  `json:"organizer"`
	Category        entities.EventCategory  `json:"category"`
	Space           EventSpace              `json:"space"`
	StartsAt        time.Time               `json:"starts_at"`
	EndsAt          time.Time               `json:"ends_at"`
	Capacity        int                     `json:"capacity"`
	Description     null.String             `json:"description,omitempty"`
	Status          entities.EventStatus    `json:"status,omitempty"`
	RegisteredCount int                     `json:"registered_count,omitempty"`
	SeatsLeft       *int                    `json:"seats_left,omitempty"`
	IsRegistered    *bool                   `json:"is_registered,omitempty"`
	IsSubscribed    *bool                   `json:"is_subscribed,omitempty"`
	State           entities.EventUserState `json:"state,omitempty"`
}

type EventList struct {
	WeekTotal int     `json:"week_total"`
	Featured  *Event  `json:"featured"`
	Items     []Event `json:"items"`
}

func NewEvent(in entities.Event, spaceName string) Event {
	return Event{
		EventID:     in.EventID,
		Title:       in.Title,
		Organizer:   in.Organizer,
		Category:    in.Category,
		Space:       EventSpace{SpaceID: in.SpaceID, Name: spaceName},
		StartsAt:    in.StartsAt,
		EndsAt:      in.EndsAt,
		Capacity:    in.Capacity,
		Description: in.Description,
		Status:      in.Status,
	}
}

func NewEventView(in entities.EventView, spaceName string, now time.Time) Event {
	seatsLeft := in.SeatsLeft()
	isRegistered := in.IsRegistered
	isSubscribed := in.IsSubscribed
	return Event{
		EventID:         in.EventID,
		Title:           in.Title,
		Organizer:       in.Organizer,
		Category:        in.Category,
		Space:           EventSpace{SpaceID: in.SpaceID, Name: spaceName},
		StartsAt:        in.StartsAt,
		EndsAt:          in.EndsAt,
		Capacity:        in.Capacity,
		Description:     in.Description,
		RegisteredCount: in.RegisteredCount,
		SeatsLeft:       &seatsLeft,
		IsRegistered:    &isRegistered,
		IsSubscribed:    &isSubscribed,
		State:           in.UserState(now),
	}
}
