package service

import (
	"context"
	"testing"
	"time"

	"github.com/guregu/null/v6"
	"github.com/kr0uch/studspace/internal/core/models/dto"
	"github.com/kr0uch/studspace/internal/core/models/entities"
	"github.com/kr0uch/studspace/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestBookingService_CreateBooking(t *testing.T) {
	ctx := context.Background()
	actor := testStudent()
	hall := testHall()
	now := time.Date(2026, 11, 1, 10, 0, 0, 0, time.UTC)
	starts, ends := tomorrowSlot(now)
	created := &entities.Booking{BookingID: 10, UserID: actor.UserID, SpaceID: hall.SpaceID, Title: "Учёба", StartsAt: starts, EndsAt: ends, Purpose: "учёба", PeopleCount: 2}

	tests := []struct {
		name          string
		people        int
		mockBehavior  func(spaces *mockSpaceRepository, bookings *mockBookingRepository)
		expectedError error
	}{
		{
			name:   "success",
			people: 2,
			mockBehavior: func(spaces *mockSpaceRepository, bookings *mockBookingRepository) {
				spaces.On("GetSpaceByID", ctx, hall.SpaceID).Return(&hall, nil)
				bookings.On("CountPendingBookingsByUser", ctx, actor.UserID, now).Return(1, nil)
				bookings.On("ListOverlappingBookings", ctx, hall.SpaceID, starts, ends, (*int64)(nil)).Return([]entities.Booking{}, nil)
				bookings.On("CreateBooking", ctx, mock.AnythingOfType("entities.Booking")).Return(created, nil)
			},
		},
		{
			name:   "failure (too many pending)",
			people: 2,
			mockBehavior: func(spaces *mockSpaceRepository, bookings *mockBookingRepository) {
				spaces.On("GetSpaceByID", ctx, hall.SpaceID).Return(&hall, nil)
				bookings.On("CountPendingBookingsByUser", ctx, actor.UserID, now).Return(3, nil)
			},
			expectedError: errors.ErrTooManyPending,
		},
		{
			name:   "success (pending overlap)",
			people: 2,
			mockBehavior: func(spaces *mockSpaceRepository, bookings *mockBookingRepository) {
				spaces.On("GetSpaceByID", ctx, hall.SpaceID).Return(&hall, nil)
				bookings.On("CountPendingBookingsByUser", ctx, actor.UserID, now).Return(0, nil)
				bookings.On("ListOverlappingBookings", ctx, hall.SpaceID, starts, ends, (*int64)(nil)).Return([]entities.Booking{{
					Status: entities.BookingStatusPending,
				}}, nil)
				bookings.On("CreateBooking", ctx, mock.AnythingOfType("entities.Booking")).Return(created, nil)
			},
		},
		{
			name:   "failure (over capacity)",
			people: hall.Capacity + 1,
			mockBehavior: func(spaces *mockSpaceRepository, bookings *mockBookingRepository) {
				spaces.On("GetSpaceByID", ctx, hall.SpaceID).Return(&hall, nil)
			},
			expectedError: errors.ErrValidation,
		},
		{
			name:   "failure (slot conflict)",
			people: 2,
			mockBehavior: func(spaces *mockSpaceRepository, bookings *mockBookingRepository) {
				spaces.On("GetSpaceByID", ctx, hall.SpaceID).Return(&hall, nil)
				bookings.On("CountPendingBookingsByUser", ctx, actor.UserID, now).Return(0, nil)
				bookings.On("ListOverlappingBookings", ctx, hall.SpaceID, starts, ends, (*int64)(nil)).Return([]entities.Booking{{
					Status: entities.BookingStatusApproved,
				}}, nil)
			},
			expectedError: errors.ErrSlotConflict,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			spaces := new(mockSpaceRepository)
			bookings := new(mockBookingRepository)
			tt.mockBehavior(spaces, bookings)

			got, err := NewBookingService(bookings, spaces, new(mockUserRepository), new(mockNotifyRepository)).
				CreateBooking(ctx, actor, dto.CreateBooking{
					SpaceID: hall.SpaceID, StartsAt: starts, EndsAt: ends, Title: "Учёба", Purpose: "учёба", PeopleCount: tt.people,
				}, now)
			if tt.expectedError != nil {
				assert.ErrorIs(t, err, tt.expectedError)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, created, got)
			}
			spaces.AssertExpectations(t)
			bookings.AssertExpectations(t)
		})
	}
}

func TestBookingService_GetBookingByID(t *testing.T) {
	ctx := context.Background()
	booking := &entities.Booking{BookingID: 10, UserID: 1}

	tests := []struct {
		name          string
		actor         entities.User
		mockBehavior  func(bookings *mockBookingRepository)
		expectedError error
	}{
		{
			name:  "success (owner)",
			actor: testStudent(),
			mockBehavior: func(bookings *mockBookingRepository) {
				bookings.On("GetBookingByID", ctx, int64(10)).Return(booking, nil)
			},
		},
		{
			name:  "success (admin)",
			actor: testAdmin(),
			mockBehavior: func(bookings *mockBookingRepository) {
				bookings.On("GetBookingByID", ctx, int64(10)).Return(booking, nil)
			},
		},
		{
			name:  "failure (foreign)",
			actor: entities.User{UserID: 99, Role: entities.UserRoleStudent},
			mockBehavior: func(bookings *mockBookingRepository) {
				bookings.On("GetBookingByID", ctx, int64(10)).Return(booking, nil)
			},
			expectedError: errors.ErrForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bookings := new(mockBookingRepository)
			tt.mockBehavior(bookings)

			got, err := NewBookingService(bookings, new(mockSpaceRepository), new(mockUserRepository), new(mockNotifyRepository)).
				GetBookingByID(ctx, tt.actor, 10)
			if tt.expectedError != nil {
				assert.ErrorIs(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, booking, got)
			}
			bookings.AssertExpectations(t)
		})
	}
}

func TestBookingService_ApproveBooking(t *testing.T) {
	ctx := context.Background()
	now := time.Date(2026, 11, 1, 10, 0, 0, 0, time.UTC)
	pending := &entities.Booking{BookingID: 10, UserID: 1, SpaceID: 1, StartsAt: now.Add(time.Hour), EndsAt: now.Add(2 * time.Hour), Status: entities.BookingStatusPending}
	approved := *pending
	approved.Status = entities.BookingStatusApproved
	hall := testHall()
	owner := testStudent()
	exclude := pending.BookingID

	t.Run("failure (not admin)", func(t *testing.T) {
		_, err := NewBookingService(new(mockBookingRepository), new(mockSpaceRepository), new(mockUserRepository), new(mockNotifyRepository)).
			ApproveBooking(ctx, testStudent(), 10, now)
		assert.ErrorIs(t, err, errors.ErrForbidden)
	})

	t.Run("success", func(t *testing.T) {
		bookings := new(mockBookingRepository)
		spaces := new(mockSpaceRepository)
		users := new(mockUserRepository)
		notify := new(mockNotifyRepository)

		bookings.On("GetBookingByID", ctx, int64(10)).Return(pending, nil)
		bookings.On("ListOverlappingBookings", ctx, pending.SpaceID, pending.StartsAt, pending.EndsAt, &exclude).Return([]entities.Booking{}, nil)
		bookings.On("ApproveBooking", ctx, int64(10), testAdmin().UserID, now).Return(&approved, nil)
		spaces.On("GetSpaceByID", ctx, pending.SpaceID).Return(&hall, nil)
		users.On("GetUserByID", ctx, owner.UserID).Return(&owner, nil)
		notify.On("Notify", mock.Anything, []int64{owner.VKUserID}, mock.AnythingOfType("string")).Return([]int64{}, nil)

		got, err := NewBookingService(bookings, spaces, users, notify).ApproveBooking(ctx, testAdmin(), 10, now)
		assert.NoError(t, err)
		assert.Equal(t, &approved, got)
		bookings.AssertExpectations(t)
		notify.AssertExpectations(t)
	})

	t.Run("failure (started)", func(t *testing.T) {
		started := *pending
		started.StartsAt = now
		bookings := new(mockBookingRepository)
		bookings.On("GetBookingByID", ctx, int64(10)).Return(&started, nil)

		_, err := NewBookingService(bookings, new(mockSpaceRepository), new(mockUserRepository), new(mockNotifyRepository)).
			ApproveBooking(ctx, testAdmin(), 10, now)
		assert.ErrorIs(t, err, errors.ErrInvalidTransition)
		bookings.AssertExpectations(t)
	})

	t.Run("success (from reschedule)", func(t *testing.T) {
		moved := *pending
		moved.RescheduledFromID = null.IntFrom(7)
		sourceID := int64(7)
		bookings := new(mockBookingRepository)
		spaces := new(mockSpaceRepository)
		users := new(mockUserRepository)
		notify := new(mockNotifyRepository)

		bookings.On("GetBookingByID", ctx, int64(10)).Return(&moved, nil)
		bookings.On("ListOverlappingBookings", ctx, moved.SpaceID, moved.StartsAt, moved.EndsAt, &sourceID).Return([]entities.Booking{}, nil)
		bookings.On("ApproveBooking", ctx, int64(10), testAdmin().UserID, now).Return(&approved, nil)
		spaces.On("GetSpaceByID", ctx, moved.SpaceID).Return(&hall, nil)
		users.On("GetUserByID", ctx, owner.UserID).Return(&owner, nil)
		notify.On("Notify", mock.Anything, []int64{owner.VKUserID}, mock.AnythingOfType("string")).Return([]int64{}, nil)

		got, err := NewBookingService(bookings, spaces, users, notify).ApproveBooking(ctx, testAdmin(), 10, now)
		assert.NoError(t, err)
		assert.Equal(t, &approved, got)
		bookings.AssertExpectations(t)
	})

	t.Run("success (notifications off)", func(t *testing.T) {
		quiet := owner
		quiet.NotificationsAllowed = false
		bookings := new(mockBookingRepository)
		spaces := new(mockSpaceRepository)
		users := new(mockUserRepository)
		notify := new(mockNotifyRepository)

		bookings.On("GetBookingByID", ctx, int64(10)).Return(pending, nil)
		bookings.On("ListOverlappingBookings", ctx, pending.SpaceID, pending.StartsAt, pending.EndsAt, &exclude).Return([]entities.Booking{}, nil)
		bookings.On("ApproveBooking", ctx, int64(10), testAdmin().UserID, now).Return(&approved, nil)
		spaces.On("GetSpaceByID", ctx, pending.SpaceID).Return(&hall, nil)
		users.On("GetUserByID", ctx, quiet.UserID).Return(&quiet, nil)

		got, err := NewBookingService(bookings, spaces, users, notify).ApproveBooking(ctx, testAdmin(), 10, now)
		assert.NoError(t, err)
		assert.Equal(t, &approved, got)
		notify.AssertNotCalled(t, "Notify", mock.Anything, mock.Anything, mock.Anything)
	})
}

func TestBookingService_RescheduleBooking(t *testing.T) {
	ctx := context.Background()
	now := time.Date(2026, 11, 1, 10, 0, 0, 0, time.UTC)
	starts, ends := tomorrowSlot(now)
	actor := testStudent()
	hall := testHall()
	current := &entities.Booking{
		BookingID: 10, UserID: actor.UserID, SpaceID: hall.SpaceID,
		Title: "Учёба", Purpose: "учёба", PeopleCount: 2,
		StartsAt: starts, EndsAt: ends, Status: entities.BookingStatusApproved,
	}
	nextStart, nextEnd := starts.Add(time.Hour), ends.Add(time.Hour)
	created := *current
	created.BookingID = 11
	created.Status = entities.BookingStatusPending
	created.StartsAt = nextStart
	created.EndsAt = nextEnd
	created.RescheduledFromID = null.IntFrom(current.BookingID)

	t.Run("success", func(t *testing.T) {
		bookings := new(mockBookingRepository)
		spaces := new(mockSpaceRepository)
		bookings.On("GetBookingByID", ctx, current.BookingID).Return(current, nil)
		spaces.On("GetSpaceByID", ctx, hall.SpaceID).Return(&hall, nil)
		bookings.On("CountPendingBookingsByUser", ctx, actor.UserID, now).Return(1, nil)
		bookings.On("HasPendingReschedule", ctx, current.BookingID).Return(false, nil)
		bookings.On("ListOverlappingBookings", ctx, hall.SpaceID, nextStart, nextEnd, &current.BookingID).Return([]entities.Booking{}, nil)
		bookings.On("RescheduleBooking", ctx, current.BookingID, actor.UserID, now, hall.SpaceID, nextStart, nextEnd).Return(&created, nil)

		got, err := NewBookingService(bookings, spaces, new(mockUserRepository), new(mockNotifyRepository)).
			RescheduleBooking(ctx, actor.UserID, dto.RescheduleBooking{
				BookingID: current.BookingID, SpaceID: hall.SpaceID, StartsAt: nextStart, EndsAt: nextEnd,
			}, now)
		assert.NoError(t, err)
		assert.Equal(t, &created, got)
		bookings.AssertExpectations(t)
	})

	t.Run("failure (too many pending)", func(t *testing.T) {
		bookings := new(mockBookingRepository)
		spaces := new(mockSpaceRepository)
		bookings.On("GetBookingByID", ctx, current.BookingID).Return(current, nil)
		spaces.On("GetSpaceByID", ctx, hall.SpaceID).Return(&hall, nil)
		bookings.On("CountPendingBookingsByUser", ctx, actor.UserID, now).Return(3, nil)

		_, err := NewBookingService(bookings, spaces, new(mockUserRepository), new(mockNotifyRepository)).
			RescheduleBooking(ctx, actor.UserID, dto.RescheduleBooking{
				BookingID: current.BookingID, SpaceID: hall.SpaceID, StartsAt: nextStart, EndsAt: nextEnd,
			}, now)
		assert.ErrorIs(t, err, errors.ErrTooManyPending)
		bookings.AssertExpectations(t)
	})

	t.Run("failure (already pending)", func(t *testing.T) {
		bookings := new(mockBookingRepository)
		spaces := new(mockSpaceRepository)
		bookings.On("GetBookingByID", ctx, current.BookingID).Return(current, nil)
		spaces.On("GetSpaceByID", ctx, hall.SpaceID).Return(&hall, nil)
		bookings.On("CountPendingBookingsByUser", ctx, actor.UserID, now).Return(1, nil)
		bookings.On("HasPendingReschedule", ctx, current.BookingID).Return(true, nil)

		_, err := NewBookingService(bookings, spaces, new(mockUserRepository), new(mockNotifyRepository)).
			RescheduleBooking(ctx, actor.UserID, dto.RescheduleBooking{
				BookingID: current.BookingID, SpaceID: hall.SpaceID, StartsAt: nextStart, EndsAt: nextEnd,
			}, now)
		assert.ErrorIs(t, err, errors.ErrValidation)
		bookings.AssertNotCalled(t, "RescheduleBooking", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("failure (not approved)", func(t *testing.T) {
		pending := *current
		pending.Status = entities.BookingStatusPending
		bookings := new(mockBookingRepository)
		spaces := new(mockSpaceRepository)
		bookings.On("GetBookingByID", ctx, current.BookingID).Return(&pending, nil)
		spaces.On("GetSpaceByID", ctx, hall.SpaceID).Return(&hall, nil)

		_, err := NewBookingService(bookings, spaces, new(mockUserRepository), new(mockNotifyRepository)).
			RescheduleBooking(ctx, actor.UserID, dto.RescheduleBooking{
				BookingID: current.BookingID, SpaceID: hall.SpaceID, StartsAt: nextStart, EndsAt: nextEnd,
			}, now)
		assert.ErrorIs(t, err, errors.ErrInvalidTransition)
	})
}

func TestBookingService_WithdrawBooking(t *testing.T) {
	ctx := context.Background()
	now := time.Date(2026, 11, 1, 10, 0, 0, 0, time.UTC)
	actor := testStudent()

	t.Run("failure (foreign)", func(t *testing.T) {
		bookings := new(mockBookingRepository)
		bookings.On("GetBookingByID", ctx, int64(10)).Return(&entities.Booking{BookingID: 10, UserID: 99}, nil)

		_, err := NewBookingService(bookings, new(mockSpaceRepository), new(mockUserRepository), new(mockNotifyRepository)).
			WithdrawBooking(ctx, actor.UserID, 10, now)
		assert.ErrorIs(t, err, errors.ErrForbidden)
		bookings.AssertNotCalled(t, "WithdrawBooking", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	})
}

func TestBookingService_RejectBooking(t *testing.T) {
	ctx := context.Background()
	now := time.Date(2026, 11, 1, 10, 0, 0, 0, time.UTC)

	t.Run("failure (not admin)", func(t *testing.T) {
		_, err := NewBookingService(new(mockBookingRepository), new(mockSpaceRepository), new(mockUserRepository), new(mockNotifyRepository)).
			RejectBooking(ctx, testStudent(), 10, entities.RejectReason{Text: "нет"}, now)
		assert.ErrorIs(t, err, errors.ErrForbidden)
	})

	t.Run("failure (empty reason)", func(t *testing.T) {
		_, err := NewBookingService(new(mockBookingRepository), new(mockSpaceRepository), new(mockUserRepository), new(mockNotifyRepository)).
			RejectBooking(ctx, testAdmin(), 10, entities.RejectReason{Text: "  "}, now)
		assert.ErrorIs(t, err, errors.ErrValidation)
	})
}
