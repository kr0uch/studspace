package service

import (
	"context"
	"strings"
	"testing"

	"github.com/guregu/null/v6"
	"github.com/kr0uch/studspace/internal/core/models/dto"
	"github.com/kr0uch/studspace/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestUserService_UpsertUserByVK(t *testing.T) {
	ctx := context.Background()
	user := testStudent()

	tests := []struct {
		name          string
		in            dto.UpsertUserByVK
		mockBehavior  func(repo *mockUserRepository)
		expectedError error
	}{
		{
			name: "success",
			in:   dto.UpsertUserByVK{VKUserID: 100, FirstName: " Иван ", LastName: " Иванов "},
			mockBehavior: func(repo *mockUserRepository) {
				repo.On("UpsertUserByVKID", ctx, int64(100), "Иван", "Иванов", null.String{}).
					Return(&user, nil)
			},
		},
		{
			name:          "failure (validation)",
			in:            dto.UpsertUserByVK{VKUserID: 0, FirstName: "", LastName: ""},
			mockBehavior:  func(repo *mockUserRepository) {},
			expectedError: errors.ErrValidation,
		},
		{
			name: "failure (db)",
			in:   dto.UpsertUserByVK{VKUserID: 100, FirstName: "Иван", LastName: "Иванов"},
			mockBehavior: func(repo *mockUserRepository) {
				repo.On("UpsertUserByVKID", ctx, int64(100), "Иван", "Иванов", null.String{}).
					Return(nil, errors.ErrInternal)
			},
			expectedError: errors.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(mockUserRepository)
			tt.mockBehavior(repo)

			got, err := NewUserService(repo).UpsertUserByVK(ctx, tt.in)
			if tt.expectedError != nil {
				assert.ErrorIs(t, err, tt.expectedError)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, &user, got)
			}
			repo.AssertExpectations(t)
		})
	}
}

func TestUserService_GetUserByID(t *testing.T) {
	ctx := context.Background()
	user := testStudent()

	tests := []struct {
		name          string
		mockBehavior  func(repo *mockUserRepository)
		expectedError error
	}{
		{
			name: "success",
			mockBehavior: func(repo *mockUserRepository) {
				repo.On("GetUserByID", ctx, int64(1)).Return(&user, nil)
			},
		},
		{
			name: "failure (not found)",
			mockBehavior: func(repo *mockUserRepository) {
				repo.On("GetUserByID", ctx, int64(1)).Return(nil, errors.ErrNotFound)
			},
			expectedError: errors.ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(mockUserRepository)
			tt.mockBehavior(repo)

			got, err := NewUserService(repo).GetUserByID(ctx, 1)
			if tt.expectedError != nil {
				assert.ErrorIs(t, err, tt.expectedError)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, &user, got)
			}
			repo.AssertExpectations(t)
		})
	}
}

func TestUserService_GetUserByVKUserID(t *testing.T) {
	ctx := context.Background()
	user := testStudent()
	repo := new(mockUserRepository)
	repo.On("GetUserByVKUserID", ctx, int64(100)).Return(&user, nil)

	got, err := NewUserService(repo).GetUserByVKUserID(ctx, 100)
	assert.NoError(t, err)
	assert.Equal(t, &user, got)
	repo.AssertExpectations(t)
}

func TestUserService_UpdateUserNotifications(t *testing.T) {
	ctx := context.Background()
	user := testStudent()
	user.NotificationsAllowed = false
	repo := new(mockUserRepository)
	repo.On("UpdateUserNotifications", ctx, int64(1), false).Return(&user, nil)

	got, err := NewUserService(repo).UpdateUserNotifications(ctx, dto.UpdateUserNotifications{UserID: 1, Allowed: false})
	assert.NoError(t, err)
	assert.Equal(t, &user, got)
	repo.AssertExpectations(t)
}

func TestUserService_UpdateUser(t *testing.T) {
	ctx := context.Background()
	user := testStudent()
	user.FirstName = "Дмитрий"
	user.LastName = "Казак"
	user.Patronymic = null.StringFrom("Алексеевич")
	user.Occupation = null.StringFrom("Студент, 25ВВБ4")

	t.Run("success", func(t *testing.T) {
		repo := new(mockUserRepository)
		repo.On("UpdateUser", ctx, int64(1), null.StringFrom("Алексеевич"), null.StringFrom("Студент, 25ВВБ4")).
			Return(&user, nil)

		got, err := NewUserService(repo).UpdateUser(ctx, dto.UpdateUser{
			UserID:     1,
			Patronymic: null.StringFrom(" Алексеевич "),
			Occupation: null.StringFrom(" Студент, 25ВВБ4 "),
		})
		assert.NoError(t, err)
		assert.Equal(t, &user, got)
		repo.AssertExpectations(t)
	})

	t.Run("failure (validation)", func(t *testing.T) {
		got, err := NewUserService(new(mockUserRepository)).UpdateUser(ctx, dto.UpdateUser{
			UserID:     1,
			Occupation: null.StringFrom(strings.Repeat("а", 151)),
		})
		assert.ErrorIs(t, err, errors.ErrValidation)
		assert.Nil(t, got)
	})
}
