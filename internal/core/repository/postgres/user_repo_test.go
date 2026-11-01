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

func userColumns() []string {
	return []string{
		"user_id", "vk_user_id", "first_name", "last_name", "photo_url",
		"role", "notifications_allowed", "created_at",
	}
}

func TestUserRepository_GetUserByID(t *testing.T) {
	ctx := context.Background()
	created := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)

	t.Run("success", func(t *testing.T) {
		db, mock := newSQLMock(t)
		mock.ExpectQuery(matchQuery(qUserGetByID)).
			WithArgs(int64(1)).
			WillReturnRows(sqlmock.NewRows(userColumns()).
				AddRow(int64(1), int64(100), "Иван", "Иванов", nil, "student", true, created))

		user, err := NewUserRepository(db).GetUserByID(ctx, 1)
		require.NoError(t, err)
		assert.Equal(t, int64(1), user.UserID)
		assert.Equal(t, int64(100), user.VKUserID)
		assert.Equal(t, "Иван", user.FirstName)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("failure (not found)", func(t *testing.T) {
		db, mock := newSQLMock(t)
		mock.ExpectQuery(matchQuery(qUserGetByID)).
			WithArgs(int64(1)).
			WillReturnError(sql.ErrNoRows)

		user, err := NewUserRepository(db).GetUserByID(ctx, 1)
		assert.ErrorIs(t, err, errors.ErrNotFound)
		assert.Nil(t, user)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestUserRepository_ListUsersByIDs(t *testing.T) {
	ctx := context.Background()

	t.Run("empty", func(t *testing.T) {
		users, err := NewUserRepository(nil).ListUsersByIDs(ctx, nil)
		assert.NoError(t, err)
		assert.Nil(t, users)
	})

	t.Run("success", func(t *testing.T) {
		db, mock := newSQLMock(t)
		created := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
		mock.ExpectQuery(matchQuery(qUserListByIDs)).
			WithArgs("{1,2}").
			WillReturnRows(sqlmock.NewRows(userColumns()).
				AddRow(int64(1), int64(100), "Иван", "Иванов", nil, "student", true, created))

		users, err := NewUserRepository(db).ListUsersByIDs(ctx, []int64{1, 2})
		require.NoError(t, err)
		assert.Len(t, users, 1)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
