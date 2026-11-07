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

func (m *mockUserRepository) UpsertUserByVKID(ctx context.Context, vkUserID int64, firstName, lastName string, photoURL null.String, role entities.UserRole) (*entities.User, error) {
	args := m.Called(ctx, vkUserID, firstName, lastName, photoURL, role)
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

func (m *mockUserRepository) UpdateUser(ctx context.Context, userID int64, patronymic, occupation null.String) (*entities.User, error) {
	args := m.Called(ctx, userID, patronymic, occupation)
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

func (m *mockBookingRepository) ListAdminBookings(ctx context.Context, tab entities.AdminBookingTab, now time.Time, cursorAt *time.Time, cursorID int64, limit int) ([]entities.Booking, error) {
	args := m.Called(ctx, tab, now, cursorAt, cursorID, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entities.Booking), args.Error(1)
}

func (m *mockBookingRepository) HasPendingReschedule(ctx context.Context, bookingID int64) (bool, error) {
	args := m.Called(ctx, bookingID)
	return args.Bool(0), args.Error(1)
}

func (m *mockBookingRepository) CountPendingBookingsByUser(ctx context.Context, userID int64, now time.Time) (int, error) {
	args := m.Called(ctx, userID, now)
	return args.Int(0), args.Error(1)
}

func (m *mockBookingRepository) CountBookingsByStatus(ctx context.Context, status entities.BookingStatus, now time.Time) (int, error) {
	args := m.Called(ctx, status, now)
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
