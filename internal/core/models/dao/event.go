package dao

import (
	"time"

	"github.com/guregu/null/v6"

	"github.com/kr0uch/studspace/internal/core/models/entities"
)

type Event struct {
	ID          int64       `db:"id"`
	Title       string      `db:"title"`
	Organizer   string      `db:"organizer"`
	Category    string      `db:"category"`
	SpaceID     int16       `db:"space_id"`
	StartsAt    time.Time   `db:"starts_at"`
	EndsAt      time.Time   `db:"ends_at"`
	Capacity    int         `db:"capacity"`
	Description null.String `db:"description"`
	Status      string      `db:"status"`
	CreatedBy   int64       `db:"created_by"`
	CreatedAt   time.Time   `db:"created_at"`
}

type EventView struct {
	Event
	RegisteredCount int  `db:"registered_count"`
	IsRegistered    bool `db:"is_registered"`
	IsSubscribed    bool `db:"is_subscribed"`
}

type EventRegistration struct {
	EventID     int64     `db:"event_id"`
	UserID      int64     `db:"user_id"`
	CreatedAt   time.Time `db:"created_at"`
	CancelledAt null.Time `db:"cancelled_at"`
}

type EventWaitlistEntry struct {
	EventID    int64     `db:"event_id"`
	UserID     int64     `db:"user_id"`
	CreatedAt  time.Time `db:"created_at"`
	NotifiedAt null.Time `db:"notified_at"`
}

func (e Event) ToEntity() entities.Event {
	return entities.Event{
		ID:          e.ID,
		Title:       e.Title,
		Organizer:   e.Organizer,
		Category:    entities.EventCategory(e.Category),
		SpaceID:     e.SpaceID,
		StartsAt:    e.StartsAt,
		EndsAt:      e.EndsAt,
		Capacity:    e.Capacity,
		Description: e.Description,
		Status:      entities.EventStatus(e.Status),
		CreatedBy:   e.CreatedBy,
		CreatedAt:   e.CreatedAt,
	}
}

func EventFromEntity(e entities.Event) Event {
	return Event{
		ID:          e.ID,
		Title:       e.Title,
		Organizer:   e.Organizer,
		Category:    string(e.Category),
		SpaceID:     e.SpaceID,
		StartsAt:    e.StartsAt,
		EndsAt:      e.EndsAt,
		Capacity:    e.Capacity,
		Description: e.Description,
		Status:      string(e.Status),
		CreatedBy:   e.CreatedBy,
		CreatedAt:   e.CreatedAt,
	}
}

func (v EventView) ToEntity() entities.EventView {
	return entities.EventView{
		Event:           v.Event.ToEntity(),
		RegisteredCount: v.RegisteredCount,
		IsRegistered:    v.IsRegistered,
		IsSubscribed:    v.IsSubscribed,
	}
}

func (r EventRegistration) ToEntity() entities.EventRegistration {
	return entities.EventRegistration{
		EventID:     r.EventID,
		UserID:      r.UserID,
		CreatedAt:   r.CreatedAt,
		CancelledAt: r.CancelledAt,
	}
}

func EventRegistrationFromEntity(e entities.EventRegistration) EventRegistration {
	return EventRegistration{
		EventID:     e.EventID,
		UserID:      e.UserID,
		CreatedAt:   e.CreatedAt,
		CancelledAt: e.CancelledAt,
	}
}

func (w EventWaitlistEntry) ToEntity() entities.EventWaitlistEntry {
	return entities.EventWaitlistEntry{
		EventID:    w.EventID,
		UserID:     w.UserID,
		CreatedAt:  w.CreatedAt,
		NotifiedAt: w.NotifiedAt,
	}
}

func EventWaitlistEntryFromEntity(e entities.EventWaitlistEntry) EventWaitlistEntry {
	return EventWaitlistEntry{
		EventID:    e.EventID,
		UserID:     e.UserID,
		CreatedAt:  e.CreatedAt,
		NotifiedAt: e.NotifiedAt,
	}
}
