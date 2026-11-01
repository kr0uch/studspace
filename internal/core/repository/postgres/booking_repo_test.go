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
		"booking_id", "user_id", "space_id", "first_name", "last_name", "patronymic", "occupation",
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
				int64(10), int64(1), int16(1), "Иван", "Иванов", nil, nil,
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
