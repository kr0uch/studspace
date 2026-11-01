package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kr0uch/studspace/internal/core/models/api"
	"github.com/kr0uch/studspace/internal/core/models/dto"
	"github.com/kr0uch/studspace/internal/core/models/entities"
	"github.com/kr0uch/studspace/pkg/web"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockHomeBookingService struct{ mock.Mock }

func (m *mockHomeBookingService) ListUserBookings(ctx context.Context, actor entities.User, in dto.ListUserBookings) ([]entities.Booking, *string, error) {
	args := m.Called(ctx, actor, in)
	if args.Get(0) == nil {
		return nil, nil, args.Error(2)
	}
	return args.Get(0).([]entities.Booking), nil, args.Error(2)
}

type mockHomeSpaceService struct{ mock.Mock }

func (m *mockHomeSpaceService) ListActiveSpaces(ctx context.Context) ([]entities.Space, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entities.Space), args.Error(1)
}

func (m *mockHomeSpaceService) GetSpaceCurrentStatus(ctx context.Context, in dto.GetSpaceCurrentStatus) (*entities.SpaceNowStatus, error) {
	args := m.Called(ctx, in)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.SpaceNowStatus), args.Error(1)
}

func TestHomeController_GetHomeHandler(t *testing.T) {
	user := testHandlerUser()
	bookings := new(mockHomeBookingService)
	spaces := new(mockHomeSpaceService)
	bookings.On("ListUserBookings", mock.Anything, *user, mock.AnythingOfType("dto.ListUserBookings")).
		Return([]entities.Booking{}, nil, nil)
	spaces.On("ListActiveSpaces", mock.Anything).Return([]entities.Space{}, nil)

	req := httptest.NewRequest(http.MethodGet, "/home", nil)
	req = req.WithContext(web.ContextWithUser(req.Context(), user))
	w := httptest.NewRecorder()
	err := NewHomeController(bookings, spaces).GetHomeHandler(w, req)
	assert.NoError(t, err)
	var got api.Home
	require.NoError(t, json.NewDecoder(w.Body).Decode(&got))
	assert.Empty(t, got.Spaces)
}
