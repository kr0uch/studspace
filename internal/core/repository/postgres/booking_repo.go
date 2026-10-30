package postgres

import (
	"context"
	"time"

	"github.com/guregu/null/v6"

	"github.com/kr0uch/studspace/internal/core/models/dao"
	"github.com/kr0uch/studspace/internal/core/models/entities"
	dbpg "github.com/kr0uch/studspace/pkg/database/postgres"
	"github.com/kr0uch/studspace/pkg/errors"
)

type BookingRepository struct {
	postgresDB *dbpg.PostgresDB
}

func NewBookingRepository(postgresDB *dbpg.PostgresDB) *BookingRepository {
	return &BookingRepository{postgresDB: postgresDB}
}

func (r *BookingRepository) CreateBooking(ctx context.Context, booking entities.Booking) (*entities.Booking, error) {
	row, err := r.getBooking(ctx, qBookingCreate,
		booking.UserID,
		booking.SpaceID,
		booking.FirstName,
		booking.LastName,
		booking.StartsAt,
		booking.EndsAt,
		booking.Purpose,
		booking.PeopleCount,
	)
	if err != nil {
		return nil, err
	}
	if err := r.insertBookingStatusLog(ctx, row.BookingID, nil, entities.BookingStatusPending, booking.UserID, null.String{}); err != nil {
		return nil, err
	}
	result := row.ToEntity()
	return &result, nil
}

func (r *BookingRepository) GetBookingByID(ctx context.Context, bookingID int64) (*entities.Booking, error) {
	row, err := r.getBooking(ctx, qBookingGetByID, bookingID)
	if err != nil {
		return nil, err
	}
	result := row.ToEntity()
	return &result, nil
}

func (r *BookingRepository) ListUserBookings(ctx context.Context, userID int64, tab entities.UserBookingTab, now time.Time, cursorAt *time.Time, cursorID int64, limit int) ([]entities.Booking, error) {
	switch tab {
	case entities.UserBookingTabHistory:
		return r.listBookings(ctx, qBookingListByUserHistory, userID, now, cursorAt, cursorID, pgLimit(limit))
	default:
		return r.listBookings(ctx, qBookingListByUserActive, userID, now, cursorAt, cursorID, pgLimit(limit))
	}
}

func (r *BookingRepository) ListAdminBookings(ctx context.Context, tab entities.AdminBookingTab, cursorAt *time.Time, cursorID int64, limit int) ([]entities.Booking, error) {
	switch tab {
	case entities.AdminBookingTabApproved:
		return r.listBookings(ctx, qBookingListAdminApproved, cursorAt, cursorID, pgLimit(limit))
	case entities.AdminBookingTabRejected:
		return r.listBookings(ctx, qBookingListAdmin, string(entities.BookingStatusRejected), cursorAt, cursorID, pgLimit(limit))
	default:
		return r.listBookings(ctx, qBookingListAdmin, string(entities.BookingStatusPending), cursorAt, cursorID, pgLimit(limit))
	}
}

func (r *BookingRepository) CountPendingBookingsByUser(ctx context.Context, userID int64) (int, error) {
	var count int
	if err := r.postgresDB.GetContext(ctx, &count, qBookingCountPendingByUser, userID); err != nil {
		return 0, wrapErr(err)
	}
	return count, nil
}

func (r *BookingRepository) CountBookingsByStatus(ctx context.Context, status entities.BookingStatus) (int, error) {
	var count int
	if err := r.postgresDB.GetContext(ctx, &count, qBookingCountByStatus, string(status)); err != nil {
		return 0, wrapErr(err)
	}
	return count, nil
}

func pgLimit(limit int) any {
	if limit < 1 {
		return nil
	}
	return limit
}

func (r *BookingRepository) ListOverlappingBookings(ctx context.Context, spaceID int16, startsAt, endsAt time.Time, excludeBookingID *int64) ([]entities.Booking, error) {
	return r.listBookings(ctx, qBookingListOverlapping, spaceID, startsAt, endsAt, excludeBookingID)
}

func (r *BookingRepository) ListBookingsInRange(ctx context.Context, from, to time.Time, spaceID *int16) ([]entities.Booking, error) {
	return r.listBookings(ctx, qBookingListInRange, from, to, spaceID)
}

func (r *BookingRepository) ListBusyBookingDays(ctx context.Context, from, to time.Time) ([]time.Time, error) {
	days := make([]time.Time, 0)
	if err := r.postgresDB.SelectContext(ctx, &days, qBookingBusyDays, from, to); err != nil {
		return nil, wrapErr(err)
	}
	return days, nil
}

func (r *BookingRepository) ApproveBooking(ctx context.Context, bookingID, actorUserID int64, now time.Time) (*entities.Booking, error) {
	return r.transitionBooking(ctx, bookingID, actorUserID, entities.BookingStatusApproved, now, null.String{}, func() (*dao.Booking, error) {
		return r.getBooking(ctx, qBookingApprove, bookingID, actorUserID)
	})
}

func (r *BookingRepository) RejectBooking(ctx context.Context, bookingID, actorUserID int64, now time.Time, reason entities.RejectReason) (*entities.Booking, error) {
	return r.transitionBooking(ctx, bookingID, actorUserID, entities.BookingStatusRejected, now, null.StringFrom(reason.Text), func() (*dao.Booking, error) {
		return r.getBooking(ctx, qBookingReject, bookingID, actorUserID, reason.Title, reason.Text)
	})
}

func (r *BookingRepository) WithdrawBooking(ctx context.Context, bookingID, actorUserID int64, now time.Time) (*entities.Booking, error) {
	return r.transitionBooking(ctx, bookingID, actorUserID, entities.BookingStatusWithdrawn, now, null.String{}, func() (*dao.Booking, error) {
		return r.getBooking(ctx, qBookingWithdraw, bookingID)
	})
}

func (r *BookingRepository) CancelBooking(ctx context.Context, bookingID, actorUserID int64, now time.Time, actor entities.CancelActor) (*entities.Booking, error) {
	return r.transitionBooking(ctx, bookingID, actorUserID, entities.BookingStatusCancelled, now, null.String{}, func() (*dao.Booking, error) {
		return r.getBooking(ctx, qBookingCancel, bookingID, string(actor))
	})
}

func (r *BookingRepository) ReturnBookingToReview(ctx context.Context, bookingID, actorUserID int64, now time.Time) (*entities.Booking, error) {
	return r.transitionBooking(ctx, bookingID, actorUserID, entities.BookingStatusPending, now, null.String{}, func() (*dao.Booking, error) {
		return r.getBooking(ctx, qBookingReturnToReview, bookingID)
	})
}

func (r *BookingRepository) RescheduleBooking(ctx context.Context, bookingID, actorUserID int64, now time.Time, spaceID int16, startsAt, endsAt time.Time) (*entities.Booking, error) {
	return r.transitionBooking(ctx, bookingID, actorUserID, entities.BookingStatusPending, now, null.String{}, func() (*dao.Booking, error) {
		return r.getBooking(ctx, qBookingReschedule, bookingID, spaceID, startsAt, endsAt)
	})
}

func (r *BookingRepository) transitionBooking(
	ctx context.Context,
	bookingID, actorUserID int64,
	next entities.BookingStatus,
	now time.Time,
	comment null.String,
	apply func() (*dao.Booking, error),
) (*entities.Booking, error) {
	current, err := r.GetBookingByID(ctx, bookingID)
	if err != nil {
		return nil, err
	}
	if !current.CurrentStatus(now).CanTransitionTo(next) {
		return nil, errors.ErrInvalidTransition
	}

	row, err := apply()
	if err != nil {
		return nil, err
	}
	if err := r.insertBookingStatusLog(ctx, bookingID, &current.Status, next, actorUserID, comment); err != nil {
		return nil, err
	}
	result := row.ToEntity()
	return &result, nil
}

func (r *BookingRepository) insertBookingStatusLog(ctx context.Context, bookingID int64, from *entities.BookingStatus, to entities.BookingStatus, actorUserID int64, comment null.String) error {
	var fromStatus any
	if from != nil {
		fromStatus = string(*from)
	}
	_, err := r.postgresDB.ExecContext(ctx, qBookingInsertStatusLog, bookingID, fromStatus, string(to), actorUserID, comment)
	return wrapErr(err)
}

func (r *BookingRepository) getBooking(ctx context.Context, query string, args ...any) (*dao.Booking, error) {
	var row dao.Booking
	if err := r.postgresDB.GetContext(ctx, &row, query, args...); err != nil {
		return nil, wrapErr(err)
	}
	return &row, nil
}

func (r *BookingRepository) listBookings(ctx context.Context, query string, args ...any) ([]entities.Booking, error) {
	var rows []dao.Booking
	if err := r.postgresDB.SelectContext(ctx, &rows, query, args...); err != nil {
		return nil, wrapErr(err)
	}
	result := make([]entities.Booking, 0, len(rows))
	for _, row := range rows {
		result = append(result, row.ToEntity())
	}
	return result, nil
}
