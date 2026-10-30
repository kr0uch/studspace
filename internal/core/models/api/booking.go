package api

import (
	"time"

	"github.com/guregu/null/v6"

	"github.com/kr0uch/studspace/internal/core/models/entities"
)

type CreateBookingRequest struct {
	SpaceID     int16     `json:"space_id"`
	StartsAt    time.Time `json:"starts_at"`
	EndsAt      time.Time `json:"ends_at"`
	Purpose     string    `json:"purpose"`
	PeopleCount int       `json:"people_count"`
}

type RescheduleBookingRequest struct {
	SpaceID  int16     `json:"space_id"`
	StartsAt time.Time `json:"starts_at"`
	EndsAt   time.Time `json:"ends_at"`
}

type RejectBookingRequest struct {
	ReasonTitle null.String `json:"reason_title"`
	ReasonText  string      `json:"reason_text"`
}

type BookingSpace struct {
	SpaceID int16  `json:"space_id"`
	Code    string `json:"code,omitempty"`
	Name    string `json:"name,omitempty"`
}

type BookingRejectReason struct {
	Title null.String `json:"title"`
	Text  string      `json:"text"`
}

type BookingDecisionBy struct {
	UserID    int64  `json:"user_id"`
	FirstName string `json:"first_name"`
}

type BookingDecision struct {
	By BookingDecisionBy `json:"by"`
	At time.Time         `json:"at"`
}

type Booking struct {
	BookingID    int64                    `json:"booking_id"`
	Status       entities.BookingStatus   `json:"status"`
	Space        BookingSpace             `json:"space"`
	FirstName    string                   `json:"first_name"`
	LastName     string                   `json:"last_name"`
	StartsAt     time.Time                `json:"starts_at"`
	EndsAt       time.Time                `json:"ends_at"`
	Purpose      string                   `json:"purpose"`
	PeopleCount  int                      `json:"people_count"`
	RejectReason *BookingRejectReason     `json:"reject_reason"`
	Decision     *BookingDecision         `json:"decision"`
	CancelledBy  *entities.CancelActor    `json:"cancelled_by"`
	CreatedAt    time.Time                `json:"created_at"`
	Actions      []entities.BookingAction `json:"actions"`
}

type BookingList struct {
	Items      []Booking `json:"items"`
	NextCursor *string   `json:"next_cursor"`
}

func NewBooking(in entities.Booking, space *entities.Space, decidedBy *entities.User, now time.Time) Booking {
	spaceView := BookingSpace{SpaceID: in.SpaceID}
	if space != nil {
		spaceView = BookingSpace{
			SpaceID: space.SpaceID,
			Code:    space.Code,
			Name:    space.Name,
		}
	}

	out := Booking{
		BookingID:   in.BookingID,
		Status:      in.CurrentStatus(now),
		Space:       spaceView,
		FirstName:   in.FirstName,
		LastName:    in.LastName,
		StartsAt:    in.StartsAt,
		EndsAt:      in.EndsAt,
		Purpose:     in.Purpose,
		PeopleCount: in.PeopleCount,
		CreatedAt:   in.CreatedAt,
		Actions:     in.StudentActions(now),
	}
	if in.RejectReason != nil {
		out.RejectReason = &BookingRejectReason{
			Title: in.RejectReason.Title,
			Text:  in.RejectReason.Text,
		}
	}
	if in.DecidedAt.Valid {
		decision := BookingDecision{
			At: in.DecidedAt.Time,
			By: BookingDecisionBy{UserID: in.DecidedBy.Int64},
		}
		if decidedBy != nil {
			decision.By.FirstName = decidedBy.FirstName
		}
		out.Decision = &decision
	}
	if in.CancelledBy.Valid {
		actor := in.CancelledBy.V
		out.CancelledBy = &actor
	}
	return out
}
