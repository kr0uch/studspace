package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/guregu/null/v6"
	"github.com/kr0uch/studspace/internal/core/models/api"
	"github.com/kr0uch/studspace/internal/core/models/dto"
	"github.com/kr0uch/studspace/internal/core/models/entities"
	"github.com/kr0uch/studspace/pkg/errors"
	"github.com/kr0uch/studspace/pkg/web"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockUserService struct{ mock.Mock }

func (m *mockUserService) UpsertUserByVK(ctx context.Context, in dto.UpsertUserByVK) (*entities.User, error) {
	args := m.Called(ctx, in)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.User), args.Error(1)
}

func (m *mockUserService) UpdateUserNotifications(ctx context.Context, in dto.UpdateUserNotifications) (*entities.User, error) {
	args := m.Called(ctx, in)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.User), args.Error(1)
}

func (m *mockUserService) UpdateUser(ctx context.Context, in dto.UpdateUser) (*entities.User, error) {
	args := m.Called(ctx, in)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.User), args.Error(1)
}

func testHandlerUser() *entities.User {
	return &entities.User{UserID: 1, VKUserID: 100, FirstName: "Иван", LastName: "Иванов", Role: entities.UserRoleStudent}
}

func TestUserController_GetMeHandler(t *testing.T) {
	user := testHandlerUser()

	tests := []struct {
		name          string
		setupContext  func(r *http.Request) *http.Request
		expectedError error
	}{
		{
			name: "success",
			setupContext: func(r *http.Request) *http.Request {
				return r.WithContext(web.ContextWithUser(r.Context(), user))
			},
		},
		{
			name:          "failure (user not in context)",
			setupContext:  func(r *http.Request) *http.Request { return r },
			expectedError: errors.ErrForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := tt.setupContext(httptest.NewRequest(http.MethodGet, "/me", nil))
			w := httptest.NewRecorder()
			err := NewUserController(new(mockUserService)).GetMeHandler(w, req)
			if tt.expectedError != nil {
				assert.ErrorIs(t, err, tt.expectedError)
				return
			}
			assert.NoError(t, err)
			var got entities.User
			require.NoError(t, json.NewDecoder(w.Body).Decode(&got))
			assert.Equal(t, user.UserID, got.UserID)
		})
	}
}

func TestUserController_SyncHandler(t *testing.T) {
	user := testHandlerUser()
	body, _ := json.Marshal(api.SyncUserRequest{FirstName: "Иван", LastName: "Иванов"})

	t.Run("success", func(t *testing.T) {
		svc := new(mockUserService)
		svc.On("UpsertUserByVK", mock.Anything, dto.UpsertUserByVK{
			VKUserID:  100,
			FirstName: "Иван",
			LastName:  "Иванов",
		}).Return(user, nil)

		req := httptest.NewRequest(http.MethodPost, "/auth/sync", bytes.NewReader(body))
		req = req.WithContext(web.ContextWithVKUserID(req.Context(), 100))
		w := httptest.NewRecorder()

		err := NewUserController(svc).SyncHandler(w, req)
		assert.NoError(t, err)
		svc.AssertExpectations(t)
	})

	t.Run("failure (bad json)", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/auth/sync", bytes.NewReader([]byte("{")))
		req = req.WithContext(web.ContextWithVKUserID(req.Context(), 100))
		err := NewUserController(new(mockUserService)).SyncHandler(httptest.NewRecorder(), req)
		assert.ErrorIs(t, err, errors.ErrBadRequest)
	})
}

func TestUserController_PatchMeHandler(t *testing.T) {
	user := testHandlerUser()
	body, _ := json.Marshal(api.PatchMeRequest{
		Patronymic: null.StringFrom("Алексеевич"),
		Occupation: null.StringFrom("Студент, 25ВВБ4"),
	})
	updated := *user
	updated.Patronymic = null.StringFrom("Алексеевич")
	updated.Occupation = null.StringFrom("Студент, 25ВВБ4")

	svc := new(mockUserService)
	svc.On("UpdateUser", mock.Anything, mock.AnythingOfType("dto.UpdateUser")).Return(&updated, nil)

	req := httptest.NewRequest(http.MethodPatch, "/me", bytes.NewReader(body))
	req = req.WithContext(web.ContextWithUser(req.Context(), user))
	w := httptest.NewRecorder()

	err := NewUserController(svc).PatchMeHandler(w, req)
	assert.NoError(t, err)
	var got entities.User
	require.NoError(t, json.NewDecoder(w.Body).Decode(&got))
	assert.Equal(t, "Алексеевич", got.Patronymic.String)
	svc.AssertExpectations(t)
}

func TestUserController_PatchNotificationsHandler(t *testing.T) {
	user := testHandlerUser()
	body, _ := json.Marshal(api.PatchNotificationsRequest{Allowed: false})
	svc := new(mockUserService)
	svc.On("UpdateUserNotifications", mock.Anything, dto.UpdateUserNotifications{UserID: 1, Allowed: false}).
		Return(user, nil)

	req := httptest.NewRequest(http.MethodPatch, "/me/notifications", bytes.NewReader(body))
	req = req.WithContext(web.ContextWithUser(req.Context(), user))
	w := httptest.NewRecorder()

	err := NewUserController(svc).PatchNotificationsHandler(w, req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, w.Code)
	svc.AssertExpectations(t)
}
