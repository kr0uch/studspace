package postgres

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/kr0uch/studspace/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func bookingColumns() []string {
	return []string{
		"booking_id", "user_id", "space_id", "first_name", "last_name", "patronymic", "occupation", "title",
		"starts_at", "ends_at", "purpose", "people_count", "status",
		"reject_reason_title", "reject_reason_text", "decided_by", "decided_at",
		"cancelled_by", "cancelled_at", "rescheduled_from_id", "created_at", "updated_at",
	}
}

func TestBookingRepository_GetBookingByID(t *testing.T) {
	ctx := context.Background()
	now := time.Date(2026, 11, 1, 12, 0, 0, 0, time.UTC)

	t.Run("success", func(t *testing.T) {
		db, mock := newSQLMock(t)
		mock.ExpectQuery(matchQuery(qBookingGetByID)).
			WithArgs(int64(10)).
			WillReturnRows(sqlmock.NewRows(bookingColumns()).AddRow(
				int64(10), int64(1), int16(1), "Иван", "Иванов", nil, nil, "Учёба",
				now, now.Add(time.Hour), "учёба", 2, "pending",
				nil, nil, nil, nil,
				nil, nil, nil, now, now,
			))

		booking, err := NewBookingRepository(db).GetBookingByID(ctx, 10)
		require.NoError(t, err)
		assert.Equal(t, int64(10), booking.BookingID)
		assert.Equal(t, "учёба", booking.Purpose)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("failure (not found)", func(t *testing.T) {
		db, mock := newSQLMock(t)
		mock.ExpectQuery(matchQuery(qBookingGetByID)).
			WithArgs(int64(10)).
			WillReturnError(sql.ErrNoRows)

		booking, err := NewBookingRepository(db).GetBookingByID(ctx, 10)
		assert.ErrorIs(t, err, errors.ErrNotFound)
		assert.Nil(t, booking)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func bookingRow(now time.Time, status string) *sqlmock.Rows {
	return sqlmock.NewRows(bookingColumns()).AddRow(
		int64(10), int64(1), int16(1), "Иван", "Иванов", nil, nil, "Учёба",
		now.Add(time.Hour), now.Add(2*time.Hour), "учёба", 2, status,
		nil, nil, nil, nil,
		nil, nil, nil, now, now,
	)
}

func TestBookingRepository_ApproveBooking(t *testing.T) {
	ctx := context.Background()
	now := time.Date(2026, 11, 1, 12, 0, 0, 0, time.UTC)

	t.Run("success", func(t *testing.T) {
		db, mock := newSQLMock(t)
		mock.ExpectBegin()
		mock.ExpectQuery(matchQuery(qBookingGetByIDForUpdate)).
			WithArgs(int64(10)).
			WillReturnRows(bookingRow(now, "pending"))
		mock.ExpectQuery(matchQuery(qBookingApprove)).
			WithArgs(int64(10), int64(2)).
			WillReturnRows(bookingRow(now, "approved"))
		mock.ExpectExec(matchQuery(qBookingInsertStatusLog)).
			WithArgs(int64(10), "pending", "approved", int64(2), sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		booking, err := NewBookingRepository(db).ApproveBooking(ctx, 10, 2, now)
		require.NoError(t, err)
		assert.Equal(t, "approved", string(booking.Status))
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("failure (invalid transition)", func(t *testing.T) {
		db, mock := newSQLMock(t)
		mock.ExpectBegin()
		mock.ExpectQuery(matchQuery(qBookingGetByIDForUpdate)).
			WithArgs(int64(10)).
			WillReturnRows(bookingRow(now, "withdrawn"))
		mock.ExpectRollback()

		booking, err := NewBookingRepository(db).ApproveBooking(ctx, 10, 2, now)
		assert.ErrorIs(t, err, errors.ErrInvalidTransition)
		assert.Nil(t, booking)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("failure (not found)", func(t *testing.T) {
		db, mock := newSQLMock(t)
		mock.ExpectBegin()
		mock.ExpectQuery(matchQuery(qBookingGetByIDForUpdate)).
			WithArgs(int64(10)).
			WillReturnError(sql.ErrNoRows)
		mock.ExpectRollback()

		booking, err := NewBookingRepository(db).ApproveBooking(ctx, 10, 2, now)
		assert.ErrorIs(t, err, errors.ErrNotFound)
		assert.Nil(t, booking)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestBookingRepository_RescheduleBooking(t *testing.T) {
	ctx := context.Background()
	now := time.Date(2026, 11, 1, 12, 0, 0, 0, time.UTC)
	starts := now.Add(2 * time.Hour)
	ends := now.Add(3 * time.Hour)

	db, mock := newSQLMock(t)
	mock.ExpectBegin()
	mock.ExpectQuery(matchQuery(qBookingGetByIDForUpdate)).
		WithArgs(int64(10)).
		WillReturnRows(bookingRow(now, "approved"))
	mock.ExpectQuery(matchQuery(qBookingReschedule)).
		WithArgs(int64(10), int16(1), starts, ends).
		WillReturnRows(bookingRow(now, "pending"))
	mock.ExpectExec(matchQuery(qBookingInsertStatusLog)).
		WithArgs(int64(10), nil, "pending", int64(1), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	booking, err := NewBookingRepository(db).RescheduleBooking(ctx, 10, 1, now, 1, starts, ends)
	require.NoError(t, err)
	assert.Equal(t, "pending", string(booking.Status))
	assert.NoError(t, mock.ExpectationsWereMet())
}
