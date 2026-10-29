package entities

import (
	"time"

	"github.com/guregu/null/v6"
)

type EventCategory string

const (
	EventCategoryStudy   EventCategory = "study"
	EventCategoryCareer  EventCategory = "career"
	EventCategoryLeisure EventCategory = "leisure"
)

type EventStatus string

const (
	EventStatusPublished EventStatus = "published"
	EventStatusCancelled EventStatus = "cancelled"
)

type Event struct {
	EventID     int64
	Title       string
	Organizer   string
	Category    EventCategory
	SpaceID     int16
	StartsAt    time.Time
	EndsAt      time.Time
	Capacity    int
	Description null.String
	Status      EventStatus
	CreatedBy   int64
	CreatedAt   time.Time
}

type EventRegistration struct {
	EventID     int64
	UserID      int64
	CreatedAt   time.Time
	CancelledAt null.Time
}

type EventWaitlistEntry struct {
	EventID    int64
	UserID     int64
	CreatedAt  time.Time
	NotifiedAt null.Time
}

type EventUserState string

const (
	EventUserStateAvailable  EventUserState = "available"
	EventUserStateFewLeft    EventUserState = "few_left"
	EventUserStateFull       EventUserState = "full"
	EventUserStateRegistered EventUserState = "registered"
	EventUserStatePast       EventUserState = "past"
)

const FewSeatsLeftBound = 5

type EventView struct {
	Event
	RegisteredCount int
	IsRegistered    bool
	IsSubscribed    bool
}

func (c EventCategory) Valid() bool {
	switch c {
	case EventCategoryStudy, EventCategoryCareer, EventCategoryLeisure:
		return true
	default:
		return false
	}
}

func (e EventView) SeatsLeft() int {
	return max(e.Capacity-e.RegisteredCount, 0)
}

func (e EventView) UserState(now time.Time) EventUserState {
	switch {
	case !e.EndsAt.After(now):
		return EventUserStatePast
	case e.IsRegistered:
		return EventUserStateRegistered
	case e.SeatsLeft() == 0:
		return EventUserStateFull
	case e.SeatsLeft() <= FewSeatsLeftBound:
		return EventUserStateFewLeft
	default:
		return EventUserStateAvailable
	}
}
