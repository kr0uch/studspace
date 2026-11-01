package service

import (
	"context"
	"time"

	"github.com/guregu/null/v6"
	"github.com/kr0uch/studspace/internal/core/models/entities"
	"github.com/kr0uch/studspace/pkg/utils"
	"github.com/stretchr/testify/mock"
)

type mockUserRepository struct{ mock.Mock }

func (m *mockUserRepository) UpsertUserByVKID(ctx context.Context, vkUserID int64, firstName, lastName string, photoURL null.String) (*entities.User, error) {
	args := m.Called(ctx, vkUserID, firstName, lastName, photoURL)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.User), args.Error(1)
}

func (m *mockUserRepository) GetUserByID(ctx context.Context, userID int64) (*entities.User, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.User), args.Error(1)
}

func (m *mockUserRepository) GetUserByVKUserID(ctx context.Context, vkUserID int64) (*entities.User, error) {
	args := m.Called(ctx, vkUserID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.User), args.Error(1)
}

func (m *mockUserRepository) UpdateUserNotifications(ctx context.Context, userID int64, allowed bool) (*entities.User, error) {
	args := m.Called(ctx, userID, allowed)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.User), args.Error(1)
}

func (m *mockUserRepository) ListUsersByIDs(ctx context.Context, userIDs []int64) ([]entities.User, error) {
	args := m.Called(ctx, userIDs)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entities.User), args.Error(1)
}

type mockSpaceRepository struct{ mock.Mock }

func (m *mockSpaceRepository) ListActiveSpaces(ctx context.Context) ([]entities.Space, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entities.Space), args.Error(1)
}

func (m *mockSpaceRepository) GetSpaceByID(ctx context.Context, spaceID int16) (*entities.Space, error) {
	args := m.Called(ctx, spaceID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Space), args.Error(1)
}

type mockBookingRepository struct{ mock.Mock }

func (m *mockBookingRepository) CreateBooking(ctx context.Context, booking entities.Booking) (*entities.Booking, error) {
	args := m.Called(ctx, booking)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Booking), args.Error(1)
}

func (m *mockBookingRepository) GetBookingByID(ctx context.Context, bookingID int64) (*entities.Booking, error) {
	args := m.Called(ctx, bookingID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Booking), args.Error(1)
}

func (m *mockBookingRepository) ListUserBookings(ctx context.Context, userID int64, tab entities.UserBookingTab, now time.Time, cursorAt *time.Time, cursorID int64, limit int) ([]entities.Booking, error) {
	args := m.Called(ctx, userID, tab, now, cursorAt, cursorID, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entities.Booking), args.Error(1)
}

func (m *mockBookingRepository) ListAdminBookings(ctx context.Context, tab entities.AdminBookingTab, cursorAt *time.Time, cursorID int64, limit int) ([]entities.Booking, error) {
	args := m.Called(ctx, tab, cursorAt, cursorID, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entities.Booking), args.Error(1)
}

func (m *mockBookingRepository) CountPendingBookingsByUser(ctx context.Context, userID int64) (int, error) {
	args := m.Called(ctx, userID)
	return args.Int(0), args.Error(1)
}

func (m *mockBookingRepository) CountBookingsByStatus(ctx context.Context, status entities.BookingStatus) (int, error) {
	args := m.Called(ctx, status)
	return args.Int(0), args.Error(1)
}

func (m *mockBookingRepository) ListOverlappingBookings(ctx context.Context, spaceID int16, startsAt, endsAt time.Time, excludeBookingID *int64) ([]entities.Booking, error) {
	args := m.Called(ctx, spaceID, startsAt, endsAt, excludeBookingID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entities.Booking), args.Error(1)
}

func (m *mockBookingRepository) ListBookingsInRange(ctx context.Context, from, to time.Time, spaceID *int16) ([]entities.Booking, error) {
	args := m.Called(ctx, from, to, spaceID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entities.Booking), args.Error(1)
}

func (m *mockBookingRepository) ListBusyBookingDays(ctx context.Context, from, to time.Time) ([]time.Time, error) {
	args := m.Called(ctx, from, to)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]time.Time), args.Error(1)
}

func (m *mockBookingRepository) ApproveBooking(ctx context.Context, bookingID, actorUserID int64, now time.Time) (*entities.Booking, error) {
	args := m.Called(ctx, bookingID, actorUserID, now)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Booking), args.Error(1)
}

func (m *mockBookingRepository) RejectBooking(ctx context.Context, bookingID, actorUserID int64, now time.Time, reason entities.RejectReason) (*entities.Booking, error) {
	args := m.Called(ctx, bookingID, actorUserID, now, reason)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Booking), args.Error(1)
}

func (m *mockBookingRepository) WithdrawBooking(ctx context.Context, bookingID, actorUserID int64, now time.Time) (*entities.Booking, error) {
	args := m.Called(ctx, bookingID, actorUserID, now)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Booking), args.Error(1)
}

func (m *mockBookingRepository) CancelBooking(ctx context.Context, bookingID, actorUserID int64, now time.Time, actor entities.CancelActor) (*entities.Booking, error) {
	args := m.Called(ctx, bookingID, actorUserID, now, actor)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Booking), args.Error(1)
}

func (m *mockBookingRepository) ReturnBookingToReview(ctx context.Context, bookingID, actorUserID int64, now time.Time) (*entities.Booking, error) {
	args := m.Called(ctx, bookingID, actorUserID, now)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Booking), args.Error(1)
}

func (m *mockBookingRepository) RescheduleBooking(ctx context.Context, bookingID, actorUserID int64, now time.Time, spaceID int16, startsAt, endsAt time.Time) (*entities.Booking, error) {
	args := m.Called(ctx, bookingID, actorUserID, now, spaceID, startsAt, endsAt)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Booking), args.Error(1)
}

type mockEventRepository struct{ mock.Mock }

func (m *mockEventRepository) GetEventView(ctx context.Context, eventID, userID int64) (*entities.EventView, error) {
	args := m.Called(ctx, eventID, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.EventView), args.Error(1)
}

func (m *mockEventRepository) ListEventViews(ctx context.Context, from, to time.Time, userID int64, category *entities.EventCategory, mine bool) ([]entities.EventView, error) {
	args := m.Called(ctx, from, to, userID, category, mine)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entities.EventView), args.Error(1)
}

func (m *mockEventRepository) ListEventsInRange(ctx context.Context, from, to time.Time, spaceID *int16) ([]entities.Event, error) {
	args := m.Called(ctx, from, to, spaceID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entities.Event), args.Error(1)
}

func (m *mockEventRepository) ListBusyEventDays(ctx context.Context, from, to time.Time) ([]time.Time, error) {
	args := m.Called(ctx, from, to)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]time.Time), args.Error(1)
}

func (m *mockEventRepository) ListOverlappingEvents(ctx context.Context, spaceID int16, startsAt, endsAt time.Time, excludeEventID *int64) ([]entities.Event, error) {
	args := m.Called(ctx, spaceID, startsAt, endsAt, excludeEventID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entities.Event), args.Error(1)
}

func (m *mockEventRepository) CreateEvent(ctx context.Context, event entities.Event) (*entities.Event, error) {
	args := m.Called(ctx, event)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Event), args.Error(1)
}

func (m *mockEventRepository) UpdateEvent(ctx context.Context, event entities.Event) (*entities.Event, error) {
	args := m.Called(ctx, event)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Event), args.Error(1)
}

func (m *mockEventRepository) CancelEvent(ctx context.Context, eventID int64) (*entities.Event, error) {
	args := m.Called(ctx, eventID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Event), args.Error(1)
}

func (m *mockEventRepository) RegisterForEvent(ctx context.Context, eventID, userID int64) (*entities.EventView, error) {
	args := m.Called(ctx, eventID, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.EventView), args.Error(1)
}

func (m *mockEventRepository) UnregisterFromEvent(ctx context.Context, eventID, userID int64) (*entities.EventView, error) {
	args := m.Called(ctx, eventID, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.EventView), args.Error(1)
}

func (m *mockEventRepository) AddEventWaitlist(ctx context.Context, eventID, userID int64) error {
	args := m.Called(ctx, eventID, userID)
	return args.Error(0)
}

func (m *mockEventRepository) RemoveEventWaitlist(ctx context.Context, eventID, userID int64) error {
	args := m.Called(ctx, eventID, userID)
	return args.Error(0)
}

func (m *mockEventRepository) ClaimEventWaitlist(ctx context.Context, eventID int64) ([]entities.EventWaitlistEntry, error) {
	args := m.Called(ctx, eventID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entities.EventWaitlistEntry), args.Error(1)
}

type mockNotifyRepository struct{ mock.Mock }

func (m *mockNotifyRepository) Notify(ctx context.Context, vkUserIDs []int64, message string) ([]int64, error) {
	args := m.Called(ctx, vkUserIDs, message)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]int64), args.Error(1)
}

func testStudent() entities.User {
	return entities.User{UserID: 1, VKUserID: 100, FirstName: "Иван", LastName: "Иванов", Role: entities.UserRoleStudent, NotificationsAllowed: true}
}

func testAdmin() entities.User {
	return entities.User{UserID: 2, VKUserID: 200, FirstName: "Админ", LastName: "Админов", Role: entities.UserRoleAdmin}
}

func testHall() entities.Space {
	return entities.Space{
		SpaceID:     1,
		Name:        "Переговорная",
		Capacity:    10,
		IsActive:    true,
		OpenTime:    8 * time.Hour,
		CloseTime:   22 * time.Hour,
		SlotStepMin: 30,
	}
}

func tomorrowSlot(now time.Time) (time.Time, time.Time) {
	day := now.In(utils.Moscow).AddDate(0, 0, 1)
	return utils.AtClock(day, 12*time.Hour), utils.AtClock(day, 13*time.Hour)
}
