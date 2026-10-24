package dao

import (
	"time"

	"github.com/guregu/null/v6"

	"github.com/kr0uch/studspace/internal/core/models/entities"
)

type Booking struct {
	ID                int64       `db:"id"`
	UserID            int64       `db:"user_id"`
	SpaceID           int16       `db:"space_id"`
	FirstName         string      `db:"first_name"`
	LastName          string      `db:"last_name"`
	StartsAt          time.Time   `db:"starts_at"`
	EndsAt            time.Time   `db:"ends_at"`
	Purpose           string      `db:"purpose"`
	PeopleCount       int         `db:"people_count"`
	Status            string      `db:"status"`
	RejectReasonTitle null.String `db:"reject_reason_title"`
	RejectReasonText  null.String `db:"reject_reason_text"`
	DecidedBy         null.Int    `db:"decided_by"`
	DecidedAt         null.Time   `db:"decided_at"`
	CancelledBy       null.String `db:"cancelled_by"`
	CancelledAt       null.Time   `db:"cancelled_at"`
	RescheduledFromID null.Int    `db:"rescheduled_from_id"`
	CreatedAt         time.Time   `db:"created_at"`
	UpdatedAt         time.Time   `db:"updated_at"`
}

type BookingStatusLog struct {
	ID          int64       `db:"id"`
	BookingID   int64       `db:"booking_id"`
	FromStatus  null.String `db:"from_status"`
	ToStatus    string      `db:"to_status"`
	ActorUserID int64       `db:"actor_user_id"`
	Comment     null.String `db:"comment"`
	CreatedAt   time.Time   `db:"created_at"`
}

func (b Booking) ToEntity() entities.Booking {
	result := entities.Booking{
		ID:                b.ID,
		UserID:            b.UserID,
		SpaceID:           b.SpaceID,
		FirstName:         b.FirstName,
		LastName:          b.LastName,
		StartsAt:          b.StartsAt,
		EndsAt:            b.EndsAt,
		Purpose:           b.Purpose,
		PeopleCount:       b.PeopleCount,
		Status:            entities.BookingStatus(b.Status),
		DecidedBy:         b.DecidedBy,
		DecidedAt:         b.DecidedAt,
		CancelledBy:       null.NewValue(entities.CancelActor(b.CancelledBy.String), b.CancelledBy.Valid),
		CancelledAt:       b.CancelledAt,
		RescheduledFromID: b.RescheduledFromID,
		CreatedAt:         b.CreatedAt,
		UpdatedAt:         b.UpdatedAt,
	}
	if b.RejectReasonText.Valid {
		result.RejectReason = &entities.RejectReason{
			Title: b.RejectReasonTitle,
			Text:  b.RejectReasonText.String,
		}
	}
	return result
}

func BookingFromEntity(e entities.Booking) Booking {
	result := Booking{
		ID:                e.ID,
		UserID:            e.UserID,
		SpaceID:           e.SpaceID,
		FirstName:         e.FirstName,
		LastName:          e.LastName,
		StartsAt:          e.StartsAt,
		EndsAt:            e.EndsAt,
		Purpose:           e.Purpose,
		PeopleCount:       e.PeopleCount,
		Status:            string(e.Status),
		DecidedBy:         e.DecidedBy,
		DecidedAt:         e.DecidedAt,
		CancelledBy:       null.NewString(string(e.CancelledBy.V), e.CancelledBy.Valid),
		CancelledAt:       e.CancelledAt,
		RescheduledFromID: e.RescheduledFromID,
		CreatedAt:         e.CreatedAt,
		UpdatedAt:         e.UpdatedAt,
	}
	if e.RejectReason != nil {
		result.RejectReasonTitle = e.RejectReason.Title
		result.RejectReasonText = null.StringFrom(e.RejectReason.Text)
	}
	return result
}

func (l BookingStatusLog) ToEntity() entities.BookingStatusLog {
	return entities.BookingStatusLog{
		ID:          l.ID,
		BookingID:   l.BookingID,
		FromStatus:  null.NewValue(entities.BookingStatus(l.FromStatus.String), l.FromStatus.Valid),
		ToStatus:    entities.BookingStatus(l.ToStatus),
		ActorUserID: l.ActorUserID,
		Comment:     l.Comment,
		CreatedAt:   l.CreatedAt,
	}
}

func BookingStatusLogFromEntity(e entities.BookingStatusLog) BookingStatusLog {
	return BookingStatusLog{
		ID:          e.ID,
		BookingID:   e.BookingID,
		FromStatus:  null.NewString(string(e.FromStatus.V), e.FromStatus.Valid),
		ToStatus:    string(e.ToStatus),
		ActorUserID: e.ActorUserID,
		Comment:     e.Comment,
		CreatedAt:   e.CreatedAt,
	}
}
