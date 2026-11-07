package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/kr0uch/studspace/internal/core/models/api"
	"github.com/kr0uch/studspace/internal/core/models/dto"
	"github.com/kr0uch/studspace/internal/core/models/entities"
	"github.com/kr0uch/studspace/pkg/utils"
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
	now := time.Now()
	from, _ := utils.DayRange(now)
	hall := entities.Space{SpaceID: 1, Name: "Пространство А"}
	bookingStart := from.Add(10 * time.Hour)

	bookings := new(mockHomeBookingService)
	spaces := new(mockHomeSpaceService)
	bookings.On("ListUserBookings", mock.Anything, *user, mock.AnythingOfType("dto.ListUserBookings")).
		Return([]entities.Booking{{
			BookingID: 10,
			SpaceID:   1,
			Title:     "Учебная работа",
			StartsAt:  bookingStart,
			EndsAt:    bookingStart.Add(2 * time.Hour),
			Purpose:   "Подготовка к защите",
			Status:    entities.BookingStatusApproved,
		}}, nil, nil)
	spaces.On("ListActiveSpaces", mock.Anything).Return([]entities.Space{hall}, nil)
	spaces.On("GetSpaceCurrentStatus", mock.Anything, mock.AnythingOfType("dto.GetSpaceCurrentStatus")).
		Return(&entities.SpaceNowStatus{State: entities.SpaceStateFree}, nil)

	req := httptest.NewRequest(http.MethodGet, "/home", nil)
	req = req.WithContext(web.ContextWithUser(req.Context(), user))
	w := httptest.NewRecorder()
	err := NewHomeController(bookings, spaces).GetHomeHandler(w, req)
	assert.NoError(t, err)
	var got api.Home
	require.NoError(t, json.NewDecoder(w.Body).Decode(&got))
	require.Len(t, got.Today.Items, 1)
	assert.Equal(t, 1, got.Today.Total)
	assert.Equal(t, int64(10), got.Today.Items[0].BookingID)
	assert.Equal(t, "Учебная работа", got.Today.Items[0].Title)
	assert.Equal(t, "Подготовка к защите", got.Today.Items[0].Purpose)
	bookings.AssertExpectations(t)
}
