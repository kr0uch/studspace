package postgres

import (
	"context"
	"strconv"
	"strings"

	"github.com/guregu/null/v6"

	"github.com/kr0uch/studspace/internal/core/models/dao"
	"github.com/kr0uch/studspace/internal/core/models/entities"
	dbpg "github.com/kr0uch/studspace/pkg/database/postgres"
)

type UserRepository struct {
	postgresDB *dbpg.PostgresDB
}

func NewUserRepository(postgresDB *dbpg.PostgresDB) *UserRepository {
	return &UserRepository{postgresDB: postgresDB}
}

func (r *UserRepository) UpsertUserByVKID(ctx context.Context, vkUserID int64, firstName, lastName string, photoURL null.String) (*entities.User, error) {
	var row dao.User
	err := r.postgresDB.GetContext(ctx, &row, qUserUpsert, vkUserID, firstName, lastName, photoURL, string(entities.UserRoleStudent))
	if err != nil {
		return nil, wrapErr(err)
	}
	user := row.ToEntity()
	return &user, nil
}

func (r *UserRepository) GetUserByID(ctx context.Context, userID int64) (*entities.User, error) {
	return r.getUser(ctx, qUserGetByID, userID)
}

func (r *UserRepository) GetUserByVKUserID(ctx context.Context, vkUserID int64) (*entities.User, error) {
	return r.getUser(ctx, qUserGetByVKUserID, vkUserID)
}

func (r *UserRepository) UpdateUserNotifications(ctx context.Context, userID int64, allowed bool) (*entities.User, error) {
	return r.getUser(ctx, qUserUpdateNotifications, userID, allowed)
}

func (r *UserRepository) UpdateUser(ctx context.Context, userID int64, patronymic, occupation null.String) (*entities.User, error) {
	return r.getUser(ctx, qUserUpdate, userID, patronymic, occupation)
}

func (r *UserRepository) ListUsersByIDs(ctx context.Context, userIDs []int64) ([]entities.User, error) {
	if len(userIDs) == 0 {
		return nil, nil
	}
	var rows []dao.User
	if err := r.postgresDB.SelectContext(ctx, &rows, qUserListByIDs, pgInt8Array(userIDs)); err != nil {
		return nil, wrapErr(err)
	}
	result := make([]entities.User, 0, len(rows))
	for _, row := range rows {
		result = append(result, row.ToEntity())
	}
	return result, nil
}

func pgInt8Array(ids []int64) string {
	var b strings.Builder
	b.WriteByte('{')
	for i, id := range ids {
		if i > 0 {
			b.WriteByte(',')
		}
		b.WriteString(strconv.FormatInt(id, 10))
	}
	b.WriteByte('}')
	return b.String()
}

func (r *UserRepository) getUser(ctx context.Context, query string, args ...any) (*entities.User, error) {
	var row dao.User
	if err := r.postgresDB.GetContext(ctx, &row, query, args...); err != nil {
		return nil, wrapErr(err)
	}
	user := row.ToEntity()
	return &user, nil
}
