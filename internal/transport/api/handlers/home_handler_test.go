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

type mockHomeEventService struct{ mock.Mock }

func (m *mockHomeEventService) ListEventViews(ctx context.Context, userID int64, in dto.GetEventViews) ([]entities.EventView, error) {
	args := m.Called(ctx, userID, in)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entities.EventView), args.Error(1)
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
	from, to := utils.DayRange(now)
	hall := entities.Space{SpaceID: 1, Name: "Пространство А"}
	bookingStart := from.Add(10 * time.Hour)
	eventStart := from.Add(16 * time.Hour)

	bookings := new(mockHomeBookingService)
	events := new(mockHomeEventService)
	spaces := new(mockHomeSpaceService)
	bookings.On("ListUserBookings", mock.Anything, *user, mock.AnythingOfType("dto.ListUserBookings")).
		Return([]entities.Booking{{
			BookingID: 10,
			SpaceID:   1,
			StartsAt:  bookingStart,
			EndsAt:    bookingStart.Add(2 * time.Hour),
			Purpose:   "Учебная работа",
			Status:    entities.BookingStatusApproved,
		}}, nil, nil)
	events.On("ListEventViews", mock.Anything, user.UserID, mock.MatchedBy(func(in dto.GetEventViews) bool {
		return in.Mine && in.From.Equal(from) && in.To.Equal(to)
	})).Return([]entities.EventView{{
		Event: entities.Event{
			EventID:  5,
			Title:    "Мастер-класс",
			SpaceID:  1,
			StartsAt: eventStart,
			EndsAt:   eventStart.Add(90 * time.Minute),
		},
		IsRegistered: true,
	}}, nil)
	spaces.On("ListActiveSpaces", mock.Anything).Return([]entities.Space{hall}, nil)
	spaces.On("GetSpaceCurrentStatus", mock.Anything, mock.AnythingOfType("dto.GetSpaceCurrentStatus")).
		Return(&entities.SpaceNowStatus{State: entities.SpaceStateFree}, nil)

	req := httptest.NewRequest(http.MethodGet, "/home", nil)
	req = req.WithContext(web.ContextWithUser(req.Context(), user))
	w := httptest.NewRecorder()
	err := NewHomeController(bookings, events, spaces).GetHomeHandler(w, req)
	assert.NoError(t, err)
	var got api.Home
	require.NoError(t, json.NewDecoder(w.Body).Decode(&got))
	require.Len(t, got.Today.Items, 2)
	assert.Equal(t, 2, got.Today.Total)
	assert.Equal(t, api.HomeItemBooking, got.Today.Items[0].Type)
	assert.Equal(t, api.HomeItemEvent, got.Today.Items[1].Type)
	assert.Equal(t, "Мастер-класс", got.Today.Items[1].Title)
	bookings.AssertExpectations(t)
	events.AssertExpectations(t)
}
