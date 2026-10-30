package entities

import (
	"time"

	"github.com/guregu/null/v6"
)

type BookingStatus string

const (
	BookingStatusPending   BookingStatus = "pending"
	BookingStatusApproved  BookingStatus = "approved"
	BookingStatusRejected  BookingStatus = "rejected"
	BookingStatusCancelled BookingStatus = "cancelled"
	BookingStatusWithdrawn BookingStatus = "withdrawn"
	BookingStatusCompleted BookingStatus = "completed"
)

var bookingTransitions = map[BookingStatus][]BookingStatus{
	BookingStatusPending:  {BookingStatusApproved, BookingStatusRejected, BookingStatusWithdrawn},
	BookingStatusApproved: {BookingStatusCancelled, BookingStatusCompleted, BookingStatusPending},
	BookingStatusRejected: {BookingStatusPending},
}

func (s BookingStatus) CanTransitionTo(status BookingStatus) bool {
	for _, val := range bookingTransitions[s] {
		if val == status {
			return true
		}
	}
	return false
}

type CancelActor string

const (
	CancelActorUser  CancelActor = "user"
	CancelActorAdmin CancelActor = "admin"
)

type BookingAction string

const (
	BookingActionWithdraw   BookingAction = "withdraw"
	BookingActionReschedule BookingAction = "reschedule"
	BookingActionCancel     BookingAction = "cancel"
	BookingActionRebook     BookingAction = "rebook"
)

type RejectReason struct {
	Title null.String
	Text  string
}

type Booking struct {
	BookingID         int64
	UserID            int64
	SpaceID           int16
	FirstName         string
	LastName          string
	StartsAt          time.Time
	EndsAt            time.Time
	Purpose           string
	PeopleCount       int
	Status            BookingStatus
	RejectReason      *RejectReason
	DecidedBy         null.Int
	DecidedAt         null.Time
	CancelledBy       null.Value[CancelActor]
	CancelledAt       null.Time
	RescheduledFromID null.Int
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

func (b Booking) CurrentStatus(now time.Time) BookingStatus {
	if b.Status == BookingStatusApproved && !b.EndsAt.After(now) {
		return BookingStatusCompleted
	}
	return b.Status
}

func (b Booking) IsActive(now time.Time) bool {
	switch b.CurrentStatus(now) {
	case BookingStatusPending, BookingStatusApproved:
		return true
	default:
		return false
	}
}

func (b Booking) StudentActions(now time.Time) []BookingAction {
	switch b.CurrentStatus(now) {
	case BookingStatusPending:
		return []BookingAction{BookingActionWithdraw}
	case BookingStatusApproved:
		return []BookingAction{BookingActionReschedule, BookingActionCancel}
	case BookingStatusRejected, BookingStatusCancelled, BookingStatusWithdrawn, BookingStatusCompleted:
		return []BookingAction{BookingActionRebook}
	default:
		return []BookingAction{}
	}
}

type UserBookingTab string

const (
	UserBookingTabActive  UserBookingTab = "active"
	UserBookingTabHistory UserBookingTab = "history"
)

type AdminBookingTab string

const (
	AdminBookingTabNew      AdminBookingTab = "new"
	AdminBookingTabApproved AdminBookingTab = "approved"
	AdminBookingTabRejected AdminBookingTab = "rejected"
)

type BookingStatusLog struct {
	BookingStatusLogID int64
	BookingID          int64
	FromStatus         null.Value[BookingStatus]
	ToStatus           BookingStatus
	ActorUserID        int64
	Comment            null.String
	CreatedAt          time.Time
}
