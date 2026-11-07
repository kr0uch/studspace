package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/kr0uch/studspace/internal/core/models/api"
	"github.com/kr0uch/studspace/internal/core/models/entities"
	"github.com/kr0uch/studspace/pkg/errors"
	"github.com/kr0uch/studspace/pkg/utils"
	"github.com/kr0uch/studspace/pkg/web"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockScheduleBookingService struct{ mock.Mock }

func (m *mockScheduleBookingService) ListBookingsInRange(ctx context.Context, from, to time.Time) ([]entities.Booking, error) {
	args := m.Called(ctx, from, to)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entities.Booking), args.Error(1)
}

func (m *mockScheduleBookingService) ListBusyBookingDays(ctx context.Context, month time.Time) ([]time.Time, error) {
	args := m.Called(ctx, month)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]time.Time), args.Error(1)
}

type mockScheduleSpaceService struct{ mock.Mock }

func (m *mockScheduleSpaceService) ListActiveSpaces(ctx context.Context) ([]entities.Space, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entities.Space), args.Error(1)
}

func TestScheduleController_GetScheduleHandler(t *testing.T) {
	user := testHandlerUser()
	tomorrow := time.Now().In(utils.Moscow).Add(24 * time.Hour)
	day := time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), 0, 0, 0, 0, utils.Moscow)
	from, to := utils.DayRange(day)
	date := day.Format("2006-01-02")

	t.Run("failure (bad from)", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/schedule?from=&to=2026-11-01", nil)
		req = req.WithContext(web.ContextWithUser(req.Context(), user))
		err := NewScheduleController(
			new(mockScheduleSpaceService),
			new(mockScheduleBookingService),
		).GetScheduleHandler(httptest.NewRecorder(), req)
		assert.ErrorIs(t, err, errors.ErrValidation)
	})

	t.Run("shows foreign pending", func(t *testing.T) {
		spaces := new(mockScheduleSpaceService)
		bookings := new(mockScheduleBookingService)
		spaces.On("ListActiveSpaces", mock.Anything).Return([]entities.Space{{
			SpaceID: 1, Name: "A", Capacity: 10, Color: "#000",
		}}, nil)
		bookings.On("ListBookingsInRange", mock.Anything, from, to).Return([]entities.Booking{
			{
				BookingID: 1, UserID: user.UserID, SpaceID: 1, Title: "Моя",
				Status: entities.BookingStatusPending, FirstName: "Иван", LastName: "Иванов",
				StartsAt: from.Add(10 * time.Hour), EndsAt: from.Add(12 * time.Hour),
			},
			{
				BookingID: 2, UserID: 99, SpaceID: 1, Title: "Чужая",
				Status: entities.BookingStatusPending, FirstName: "Пётр", LastName: "Петров",
				StartsAt: from.Add(13 * time.Hour), EndsAt: from.Add(14 * time.Hour),
			},
			{
				BookingID: 3, UserID: 99, SpaceID: 1, Title: "Чужая ок",
				Status: entities.BookingStatusApproved, FirstName: "Пётр", LastName: "Петров",
				StartsAt: from.Add(15 * time.Hour), EndsAt: from.Add(16 * time.Hour),
			},
		}, nil)

		req := httptest.NewRequest(http.MethodGet, "/schedule?from="+date+"&to="+date, nil)
		req = req.WithContext(web.ContextWithUser(req.Context(), user))
		w := httptest.NewRecorder()
		err := NewScheduleController(spaces, bookings).GetScheduleHandler(w, req)
		require.NoError(t, err)

		var got api.Schedule
		require.NoError(t, json.NewDecoder(w.Body).Decode(&got))
		require.Len(t, got.Items, 3)
		assert.Equal(t, int64(1), got.Items[0].BookingID)
		assert.True(t, got.Items[0].IsMine)
		assert.Equal(t, "Моя", got.Items[0].Title)
		assert.Equal(t, int64(2), got.Items[1].BookingID)
		assert.False(t, got.Items[1].IsMine)
		assert.Equal(t, entities.BookingStatusPending, got.Items[1].Status)
		assert.Equal(t, int64(3), got.Items[2].BookingID)
		assert.False(t, got.Items[2].IsMine)
		assert.Equal(t, "A", got.Spaces[0].Name)
		spaces.AssertExpectations(t)
		bookings.AssertExpectations(t)
	})
}

func TestScheduleController_GetBusyDaysHandler(t *testing.T) {
	t.Run("failure (bad month)", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/calendar/busy-days?month=11", nil)
		err := NewScheduleController(
			new(mockScheduleSpaceService),
			new(mockScheduleBookingService),
		).GetBusyDaysHandler(httptest.NewRecorder(), req)
		assert.ErrorIs(t, err, errors.ErrValidation)
	})

	t.Run("unique days", func(t *testing.T) {
		month := time.Date(2026, 11, 1, 0, 0, 0, 0, utils.Moscow)
		d1 := time.Date(2026, 11, 1, 0, 0, 0, 0, utils.Moscow)
		d2 := time.Date(2026, 11, 2, 0, 0, 0, 0, utils.Moscow)
		bookings := new(mockScheduleBookingService)
		bookings.On("ListBusyBookingDays", mock.Anything, month).Return([]time.Time{d1, d2, d2}, nil)

		req := httptest.NewRequest(http.MethodGet, "/calendar/busy-days?month=2026-11", nil)
		w := httptest.NewRecorder()
		err := NewScheduleController(new(mockScheduleSpaceService), bookings).GetBusyDaysHandler(w, req)
		require.NoError(t, err)

		var got struct {
			Month string   `json:"month"`
			Days  []string `json:"days"`
		}
		require.NoError(t, json.NewDecoder(w.Body).Decode(&got))
		assert.Equal(t, "2026-11", got.Month)
		assert.Equal(t, []string{"2026-11-01", "2026-11-02"}, got.Days)
		bookings.AssertExpectations(t)
	})
}
