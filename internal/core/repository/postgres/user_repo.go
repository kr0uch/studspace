package postgres

import (
	"context"

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
	err := r.postgresDB.GetContext(ctx, &row, qUserUpsert, vkUserID, firstName, lastName, photoURL)
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

func (r *UserRepository) getUser(ctx context.Context, query string, args ...any) (*entities.User, error) {
	var row dao.User
	if err := r.postgresDB.GetContext(ctx, &row, query, args...); err != nil {
		return nil, wrapErr(err)
	}
	user := row.ToEntity()
	return &user, nil
}
