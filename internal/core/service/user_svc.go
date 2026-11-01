package service

import (
	"context"
	"strings"

	"github.com/guregu/null/v6"

	"github.com/kr0uch/studspace/internal/core/models/dto"
	"github.com/kr0uch/studspace/internal/core/models/entities"
	"github.com/kr0uch/studspace/pkg/errors"
)

type UserRepository interface {
	UpsertUserByVKID(ctx context.Context, vkUserID int64, firstName, lastName string, photoURL null.String) (*entities.User, error)
	GetUserByID(ctx context.Context, userID int64) (*entities.User, error)
	GetUserByVKUserID(ctx context.Context, vkUserID int64) (*entities.User, error)
	UpdateUserNotifications(ctx context.Context, userID int64, allowed bool) (*entities.User, error)
	ListUsersByIDs(ctx context.Context, userIDs []int64) ([]entities.User, error)
}

type UserService struct {
	users UserRepository
}

func NewUserService(users UserRepository) *UserService {
	return &UserService{users: users}
}

func (s *UserService) UpsertUserByVK(ctx context.Context, in dto.UpsertUserByVK) (*entities.User, error) {
	in.FirstName = strings.TrimSpace(in.FirstName)
	in.LastName = strings.TrimSpace(in.LastName)

	fields := map[string]string{}
	if in.VKUserID < 1 {
		fields["vk_user_id"] = "некорректный идентификатор"
	}
	if in.FirstName == "" {
		fields["first_name"] = "обязательное поле"
	}
	if in.LastName == "" {
		fields["last_name"] = "обязательное поле"
	}
	if len(fields) > 0 {
		return nil, errors.Validation(fields)
	}

	return s.users.UpsertUserByVKID(ctx, in.VKUserID, in.FirstName, in.LastName, in.PhotoURL)
}

func (s *UserService) GetUserByID(ctx context.Context, userID int64) (*entities.User, error) {
	return s.users.GetUserByID(ctx, userID)
}

func (s *UserService) GetUserByVKUserID(ctx context.Context, vkUserID int64) (*entities.User, error) {
	return s.users.GetUserByVKUserID(ctx, vkUserID)
}

func (s *UserService) UpdateUserNotifications(ctx context.Context, in dto.UpdateUserNotifications) (*entities.User, error) {
	return s.users.UpdateUserNotifications(ctx, in.UserID, in.Allowed)
}
