package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/kr0uch/studspace/internal/core/models/api"
	"github.com/kr0uch/studspace/internal/core/models/dto"
	"github.com/kr0uch/studspace/internal/core/models/entities"
	"github.com/kr0uch/studspace/pkg/errors"
	"github.com/kr0uch/studspace/pkg/utils"
	"github.com/kr0uch/studspace/pkg/web"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockBookingService struct{ mock.Mock }

func (m *mockBookingService) CreateBooking(ctx context.Context, actor entities.User, in dto.CreateBooking, now time.Time) (*entities.Booking, error) {
	args := m.Called(ctx, actor, in, now)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Booking), args.Error(1)
}

func (m *mockBookingService) GetBookingByID(ctx context.Context, actor entities.User, bookingID int64) (*entities.Booking, error) {
	args := m.Called(ctx, actor, bookingID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Booking), args.Error(1)
}

func (m *mockBookingService) ListBookingsInRange(ctx context.Context, from, to time.Time) ([]entities.Booking, error) {
	args := m.Called(ctx, from, to)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entities.Booking), args.Error(1)
}

func (m *mockBookingService) ListUserBookings(ctx context.Context, actor entities.User, in dto.ListUserBookings) ([]entities.Booking, *string, error) {
	args := m.Called(ctx, actor, in)
	var cursor *string
	if args.Get(1) != nil {
		cursor = args.Get(1).(*string)
	}
	if args.Get(0) == nil {
		return nil, cursor, args.Error(2)
	}
	return args.Get(0).([]entities.Booking), cursor, args.Error(2)
}

func (m *mockBookingService) WithdrawBooking(ctx context.Context, userID, bookingID int64, now time.Time) (*entities.Booking, error) {
	args := m.Called(ctx, userID, bookingID, now)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Booking), args.Error(1)
}

func (m *mockBookingService) CancelBookingByUser(ctx context.Context, userID, bookingID int64, now time.Time) (*entities.Booking, error) {
	args := m.Called(ctx, userID, bookingID, now)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Booking), args.Error(1)
}

func (m *mockBookingService) RescheduleBooking(ctx context.Context, userID int64, in dto.RescheduleBooking, now time.Time) (*entities.Booking, error) {
	args := m.Called(ctx, userID, in, now)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Booking), args.Error(1)
}

type mockBookingSpaceService struct{ mock.Mock }

func (m *mockBookingSpaceService) ListActiveSpaces(ctx context.Context) ([]entities.Space, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entities.Space), args.Error(1)
}

func (m *mockBookingSpaceService) GetSpaceByID(ctx context.Context, spaceID int16) (*entities.Space, error) {
	args := m.Called(ctx, spaceID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Space), args.Error(1)
}

type mockBookingUserService struct{ mock.Mock }

func (m *mockBookingUserService) GetUserByID(ctx context.Context, userID int64) (*entities.User, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.User), args.Error(1)
}

func TestBookingController_CreateBookingHandler(t *testing.T) {
	user := testHandlerUser()
	starts := time.Date(2026, 11, 2, 12, 0, 0, 0, time.UTC)
	ends := starts.Add(time.Hour)
	body, _ := json.Marshal(api.CreateBookingRequest{
		SpaceID: 1, StartsAt: starts, EndsAt: ends, Title: "Учёба", Purpose: "учёба", PeopleCount: 2,
	})
	created := &entities.Booking{BookingID: 10, UserID: 1, SpaceID: 1, Title: "Учёба", StartsAt: starts, EndsAt: ends, Purpose: "учёба", PeopleCount: 2, Status: entities.BookingStatusPending}

	t.Run("success", func(t *testing.T) {
		bookings := new(mockBookingService)
		spaces := new(mockBookingSpaceService)
		bookings.On("CreateBooking", mock.Anything, *user, dto.CreateBooking{
			SpaceID: 1, StartsAt: starts, EndsAt: ends, Title: "Учёба", Purpose: "учёба", PeopleCount: 2,
		}, mock.AnythingOfType("time.Time")).Return(created, nil)
		spaces.On("GetSpaceByID", mock.Anything, int16(1)).Return(&entities.Space{SpaceID: 1, Name: "A"}, nil)

		req := httptest.NewRequest(http.MethodPost, "/bookings", bytes.NewReader(body))
		req = req.WithContext(web.ContextWithUser(req.Context(), user))
		w := httptest.NewRecorder()
		err := NewBookingController(bookings, spaces, new(mockBookingUserService)).CreateBookingHandler(w, req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, w.Code)
		bookings.AssertExpectations(t)
	})

	t.Run("failure (slot conflict)", func(t *testing.T) {
		bookings := new(mockBookingService)
		bookings.On("CreateBooking", mock.Anything, *user, mock.Anything, mock.Anything).
			Return(nil, errors.ErrSlotConflict)

		req := httptest.NewRequest(http.MethodPost, "/bookings", bytes.NewReader(body))
		req = req.WithContext(web.ContextWithUser(req.Context(), user))
		err := NewBookingController(bookings, new(mockBookingSpaceService), new(mockBookingUserService)).CreateBookingHandler(httptest.NewRecorder(), req)
		assert.ErrorIs(t, err, errors.ErrSlotConflict)
	})
}

func TestBookingController_ListMyBookingsHandler(t *testing.T) {
	user := testHandlerUser()
	bookings := new(mockBookingService)
	bookings.On("ListUserBookings", mock.Anything, *user, mock.AnythingOfType("dto.ListUserBookings")).
		Return([]entities.Booking{}, nil, nil)

	req := httptest.NewRequest(http.MethodGet, "/bookings/my?tab=active", nil)
	req = req.WithContext(web.ContextWithUser(req.Context(), user))
	w := httptest.NewRecorder()
	err := NewBookingController(bookings, new(mockBookingSpaceService), new(mockBookingUserService)).ListMyBookingsHandler(w, req)
	assert.NoError(t, err)
	var got api.BookingList
	require.NoError(t, json.NewDecoder(w.Body).Decode(&got))
	assert.Empty(t, got.Items)
}

func TestBookingController_ListBookingsHandler(t *testing.T) {
	user := testHandlerUser()
	now := time.Now()
	from, to := utils.DayRange(now)
	date := now.In(utils.Moscow).Format("2006-01-02")
	featuredStart := now.Add(time.Hour)

	t.Run("failure (bad from)", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/bookings?from=&to=2026-11-01", nil)
		req = req.WithContext(web.ContextWithUser(req.Context(), user))
		err := NewBookingController(new(mockBookingService), new(mockBookingSpaceService), new(mockBookingUserService)).
			ListBookingsHandler(httptest.NewRecorder(), req)
		assert.ErrorIs(t, err, errors.ErrValidation)
	})

	t.Run("approved only and featured today", func(t *testing.T) {
		bookings := new(mockBookingService)
		spaces := new(mockBookingSpaceService)
		bookings.On("ListBookingsInRange", mock.Anything, from, to).Return([]entities.Booking{
			{
				BookingID: 1, SpaceID: 1, Title: "Митап", Purpose: "для команды",
				Status: entities.BookingStatusApproved, FirstName: "Иван", LastName: "Иванов",
				StartsAt: featuredStart, EndsAt: featuredStart.Add(time.Hour), PeopleCount: 4,
			},
			{
				BookingID: 2, SpaceID: 1, Title: "Черновик",
				Status:   entities.BookingStatusPending,
				StartsAt: featuredStart.Add(2 * time.Hour), EndsAt: featuredStart.Add(3 * time.Hour),
			},
		}, nil)
		spaces.On("ListActiveSpaces", mock.Anything).Return([]entities.Space{{SpaceID: 1, Name: "A"}}, nil)

		req := httptest.NewRequest(http.MethodGet, "/bookings?from="+date+"&to="+date, nil)
		req = req.WithContext(web.ContextWithUser(req.Context(), user))
		w := httptest.NewRecorder()
		err := NewBookingController(bookings, spaces, new(mockBookingUserService)).ListBookingsHandler(w, req)
		require.NoError(t, err)

		var got api.BookingFeed
		require.NoError(t, json.NewDecoder(w.Body).Decode(&got))
		require.Len(t, got.Items, 1)
		assert.Equal(t, 1, got.WeekTotal)
		assert.Equal(t, "Митап", got.Items[0].Title)
		if featuredStart.In(utils.Moscow).Format("2006-01-02") == date {
			require.NotNil(t, got.Featured)
			assert.Equal(t, int64(1), got.Featured.BookingID)
		}
		assert.Equal(t, "для команды", got.Items[0].Purpose)
		bookings.AssertExpectations(t)
		spaces.AssertExpectations(t)
	})
}

func TestBookingController_WithdrawBookingHandler(t *testing.T) {
	user := testHandlerUser()

	t.Run("failure (bad id)", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/bookings/nope/withdraw", nil)
		req.SetPathValue("id", "nope")
		req = req.WithContext(web.ContextWithUser(req.Context(), user))

		err := NewBookingController(new(mockBookingService), new(mockBookingSpaceService), new(mockBookingUserService)).
			WithdrawBookingHandler(httptest.NewRecorder(), req)
		assert.ErrorIs(t, err, errors.ErrValidation)
	})

	t.Run("success", func(t *testing.T) {
		starts := time.Date(2026, 11, 2, 12, 0, 0, 0, time.UTC)
		withdrawn := &entities.Booking{
			BookingID: 10, UserID: user.UserID, SpaceID: 1, Title: "Учёба",
			StartsAt: starts, EndsAt: starts.Add(time.Hour), Status: entities.BookingStatusWithdrawn,
		}
		bookings := new(mockBookingService)
		spaces := new(mockBookingSpaceService)
		bookings.On("WithdrawBooking", mock.Anything, user.UserID, int64(10), mock.AnythingOfType("time.Time")).Return(withdrawn, nil)
		spaces.On("GetSpaceByID", mock.Anything, int16(1)).Return(&entities.Space{SpaceID: 1, Name: "A"}, nil)

		req := httptest.NewRequest(http.MethodPost, "/bookings/10/withdraw", nil)
		req.SetPathValue("id", "10")
		req = req.WithContext(web.ContextWithUser(req.Context(), user))
		w := httptest.NewRecorder()

		err := NewBookingController(bookings, spaces, new(mockBookingUserService)).WithdrawBookingHandler(w, req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, w.Code)
		bookings.AssertExpectations(t)
	})
}

func TestBookingController_RescheduleBookingHandler(t *testing.T) {
	user := testHandlerUser()

	t.Run("failure (bad body)", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/bookings/10/reschedule", bytes.NewBufferString("{"))
		req.SetPathValue("id", "10")
		req = req.WithContext(web.ContextWithUser(req.Context(), user))

		err := NewBookingController(new(mockBookingService), new(mockBookingSpaceService), new(mockBookingUserService)).
			RescheduleBookingHandler(httptest.NewRecorder(), req)
		assert.ErrorIs(t, err, errors.ErrBadRequest)
	})
}
