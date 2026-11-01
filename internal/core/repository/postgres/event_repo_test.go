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

func TestEventRepository_GetEventView(t *testing.T) {
	ctx := context.Background()
	now := time.Date(2026, 11, 1, 12, 0, 0, 0, time.UTC)
	cols := []string{
		"event_id", "title", "organizer", "category", "space_id",
		"starts_at", "ends_at", "capacity", "description", "status",
		"created_by", "created_at", "registered_count", "is_registered", "is_subscribed",
	}

	t.Run("success", func(t *testing.T) {
		db, mock := newSQLMock(t)
		mock.ExpectQuery(matchQuery(qEventGetView)).
			WithArgs(int64(5), int64(1)).
			WillReturnRows(sqlmock.NewRows(cols).AddRow(
				int64(5), "Митап", "Студсовет", "study", int16(1),
				now, now.Add(time.Hour), 20, nil, "published",
				int64(2), now, 3, false, false,
			))

		view, err := NewEventRepository(db).GetEventView(ctx, 5, 1)
		require.NoError(t, err)
		assert.Equal(t, "Митап", view.Title)
		assert.Equal(t, 3, view.RegisteredCount)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("failure (not found)", func(t *testing.T) {
		db, mock := newSQLMock(t)
		mock.ExpectQuery(matchQuery(qEventGetView)).
			WithArgs(int64(5), int64(1)).
			WillReturnError(sql.ErrNoRows)

		view, err := NewEventRepository(db).GetEventView(ctx, 5, 1)
		assert.ErrorIs(t, err, errors.ErrNotFound)
		assert.Nil(t, view)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
