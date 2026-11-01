package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/kr0uch/studspace/internal/core/models/dto"
	"github.com/kr0uch/studspace/internal/core/models/entities"
	"github.com/kr0uch/studspace/pkg/errors"
	"github.com/kr0uch/studspace/pkg/web"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockAdminBookingService struct{ mock.Mock }

func (m *mockAdminBookingService) ListAdminBookings(ctx context.Context, actor entities.User, in dto.ListAdminBookings) ([]entities.Booking, *string, error) {
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

func (m *mockAdminBookingService) CountBookingsByStatus(ctx context.Context, actor entities.User, status entities.BookingStatus) (int, error) {
	args := m.Called(ctx, actor, status)
	return args.Int(0), args.Error(1)
}

func (m *mockAdminBookingService) GetSlotOverlaps(ctx context.Context, actor entities.User, in dto.GetSlotOverlaps) (*dto.SlotOverlaps, error) {
	args := m.Called(ctx, actor, in)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.SlotOverlaps), args.Error(1)
}

func (m *mockAdminBookingService) ApproveBooking(ctx context.Context, actor entities.User, bookingID int64, now time.Time) (*entities.Booking, error) {
	args := m.Called(ctx, actor, bookingID, now)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Booking), args.Error(1)
}

func (m *mockAdminBookingService) RejectBooking(ctx context.Context, actor entities.User, bookingID int64, reason entities.RejectReason, now time.Time) (*entities.Booking, error) {
	args := m.Called(ctx, actor, bookingID, reason, now)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Booking), args.Error(1)
}

func (m *mockAdminBookingService) CancelBookingByAdmin(ctx context.Context, actor entities.User, bookingID int64, now time.Time) (*entities.Booking, error) {
	args := m.Called(ctx, actor, bookingID, now)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Booking), args.Error(1)
}

func (m *mockAdminBookingService) ReturnBookingToReview(ctx context.Context, actor entities.User, bookingID int64, now time.Time) (*entities.Booking, error) {
	args := m.Called(ctx, actor, bookingID, now)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Booking), args.Error(1)
}

type mockAdminEventService struct{ mock.Mock }

func (m *mockAdminEventService) CreateEvent(ctx context.Context, actor entities.User, in dto.CreateEvent, now time.Time) (*entities.Event, error) {
	args := m.Called(ctx, actor, in, now)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Event), args.Error(1)
}

func (m *mockAdminEventService) UpdateEvent(ctx context.Context, actor entities.User, eventID int64, in dto.CreateEvent, now time.Time) (*entities.Event, error) {
	args := m.Called(ctx, actor, eventID, in, now)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Event), args.Error(1)
}

func (m *mockAdminEventService) CancelEvent(ctx context.Context, actor entities.User, eventID int64) (*entities.Event, error) {
	args := m.Called(ctx, actor, eventID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Event), args.Error(1)
}

type mockAdminSpaceService struct{ mock.Mock }

func (m *mockAdminSpaceService) GetSpaceByID(ctx context.Context, spaceID int16) (*entities.Space, error) {
	args := m.Called(ctx, spaceID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Space), args.Error(1)
}

type mockAdminUserService struct{ mock.Mock }

func (m *mockAdminUserService) GetUserByID(ctx context.Context, userID int64) (*entities.User, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.User), args.Error(1)
}

func TestAdminController_ApproveBookingHandler(t *testing.T) {
	admin := &entities.User{UserID: 2, Role: entities.UserRoleAdmin, FirstName: "Админ", LastName: "Админов"}
	booking := &entities.Booking{BookingID: 10, UserID: 1, SpaceID: 1, Status: entities.BookingStatusApproved}

	t.Run("success", func(t *testing.T) {
		bookings := new(mockAdminBookingService)
		spaces := new(mockAdminSpaceService)
		users := new(mockAdminUserService)
		bookings.On("ApproveBooking", mock.Anything, *admin, int64(10), mock.AnythingOfType("time.Time")).Return(booking, nil)
		spaces.On("GetSpaceByID", mock.Anything, int16(1)).Return(&entities.Space{SpaceID: 1, Name: "A"}, nil)
		users.On("GetUserByID", mock.Anything, int64(1)).Return(&entities.User{UserID: 1, FirstName: "Иван"}, nil)
		bookings.On("GetSlotOverlaps", mock.Anything, *admin, mock.AnythingOfType("dto.GetSlotOverlaps")).
			Return(&dto.SlotOverlaps{}, nil)

		req := httptest.NewRequest(http.MethodPost, "/admin/bookings/10/approve", nil)
		req.SetPathValue("id", "10")
		req = req.WithContext(web.ContextWithUser(req.Context(), admin))
		w := httptest.NewRecorder()
		err := NewAdminController(bookings, new(mockAdminEventService), spaces, users).ApproveBookingHandler(w, req)
		assert.NoError(t, err)
		bookings.AssertExpectations(t)
	})

	t.Run("failure (forbidden)", func(t *testing.T) {
		bookings := new(mockAdminBookingService)
		bookings.On("ApproveBooking", mock.Anything, mock.Anything, int64(10), mock.Anything).
			Return(nil, errors.ErrForbidden)
		req := httptest.NewRequest(http.MethodPost, "/admin/bookings/10/approve", nil)
		req.SetPathValue("id", "10")
		req = req.WithContext(web.ContextWithUser(req.Context(), testHandlerUser()))
		err := NewAdminController(bookings, new(mockAdminEventService), new(mockAdminSpaceService), new(mockAdminUserService)).
			ApproveBookingHandler(httptest.NewRecorder(), req)
		assert.ErrorIs(t, err, errors.ErrForbidden)
	})
}
