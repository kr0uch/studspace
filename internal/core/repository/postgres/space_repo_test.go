package postgres

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/kr0uch/studspace/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSpaceRepository_GetSpaceByID(t *testing.T) {
	ctx := context.Background()
	cols := []string{
		"space_id", "code", "name", "kind", "capacity", "photo_url", "color",
		"open_time", "close_time", "slot_step_min", "is_active",
	}
	amenityCols := []string{"space_id", "amenity_id", "name", "icon"}

	t.Run("success", func(t *testing.T) {
		db, mock := newSQLMock(t)
		mock.ExpectQuery(matchQuery(qSpaceGetByID)).
			WithArgs(int16(1)).
			WillReturnRows(sqlmock.NewRows(cols).AddRow(
				int16(1), "A", "Переговорная", "room", 10, nil, "#000",
				"08:00:00", "22:00:00", 30, true,
			))
		mock.ExpectQuery(matchQuery(qSpaceListAmenitiesBySpace)).
			WithArgs(int16(1)).
			WillReturnRows(sqlmock.NewRows(amenityCols).AddRow(int16(1), int16(7), "TV", "tv"))

		space, err := NewSpaceRepository(db).GetSpaceByID(ctx, 1)
		require.NoError(t, err)
		assert.Equal(t, "Переговорная", space.Name)
		assert.True(t, space.IsActive)
		assert.Len(t, space.Amenities, 1)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("failure (not found)", func(t *testing.T) {
		db, mock := newSQLMock(t)
		mock.ExpectQuery(matchQuery(qSpaceGetByID)).
			WithArgs(int16(1)).
			WillReturnError(sql.ErrNoRows)

		space, err := NewSpaceRepository(db).GetSpaceByID(ctx, 1)
		assert.ErrorIs(t, err, errors.ErrNotFound)
		assert.Nil(t, space)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
