package service

import (
	"context"
	"testing"
	"time"

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
	created := &entities.Booking{BookingID: 10, UserID: actor.UserID, SpaceID: hall.SpaceID, StartsAt: starts, EndsAt: ends, Purpose: "учёба", PeopleCount: 2}

	tests := []struct {
		name          string
		people        int
		mockBehavior  func(spaces *mockSpaceRepository, bookings *mockBookingRepository, events *mockEventRepository)
		expectedError error
	}{
		{
			name:   "success",
			people: 2,
			mockBehavior: func(spaces *mockSpaceRepository, bookings *mockBookingRepository, events *mockEventRepository) {
				spaces.On("GetSpaceByID", ctx, hall.SpaceID).Return(&hall, nil)
				bookings.On("CountPendingBookingsByUser", ctx, actor.UserID).Return(1, nil)
				bookings.On("ListOverlappingBookings", ctx, hall.SpaceID, starts, ends, (*int64)(nil)).Return([]entities.Booking{}, nil)
				events.On("ListOverlappingEvents", ctx, hall.SpaceID, starts, ends, (*int64)(nil)).Return([]entities.Event{}, nil)
				bookings.On("CreateBooking", ctx, mock.AnythingOfType("entities.Booking")).Return(created, nil)
			},
		},
		{
			name:   "failure (too many pending)",
			people: 2,
			mockBehavior: func(spaces *mockSpaceRepository, bookings *mockBookingRepository, events *mockEventRepository) {
				spaces.On("GetSpaceByID", ctx, hall.SpaceID).Return(&hall, nil)
				bookings.On("CountPendingBookingsByUser", ctx, actor.UserID).Return(3, nil)
			},
			expectedError: errors.ErrTooManyPending,
		},
		{
			name:   "failure (slot conflict)",
			people: 2,
			mockBehavior: func(spaces *mockSpaceRepository, bookings *mockBookingRepository, events *mockEventRepository) {
				spaces.On("GetSpaceByID", ctx, hall.SpaceID).Return(&hall, nil)
				bookings.On("CountPendingBookingsByUser", ctx, actor.UserID).Return(0, nil)
				bookings.On("ListOverlappingBookings", ctx, hall.SpaceID, starts, ends, (*int64)(nil)).Return([]entities.Booking{{
					Status: entities.BookingStatusApproved,
				}}, nil)
				events.On("ListOverlappingEvents", ctx, hall.SpaceID, starts, ends, (*int64)(nil)).Return([]entities.Event{}, nil)
			},
			expectedError: errors.ErrSlotConflict,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			spaces := new(mockSpaceRepository)
			bookings := new(mockBookingRepository)
			events := new(mockEventRepository)
			tt.mockBehavior(spaces, bookings, events)

			got, err := NewBookingService(bookings, events, spaces, new(mockUserRepository), new(mockNotifyRepository)).
				CreateBooking(ctx, actor, dto.CreateBooking{
					SpaceID: hall.SpaceID, StartsAt: starts, EndsAt: ends, Purpose: "учёба", PeopleCount: tt.people,
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
			events.AssertExpectations(t)
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

			got, err := NewBookingService(bookings, new(mockEventRepository), new(mockSpaceRepository), new(mockUserRepository), new(mockNotifyRepository)).
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
		_, err := NewBookingService(new(mockBookingRepository), new(mockEventRepository), new(mockSpaceRepository), new(mockUserRepository), new(mockNotifyRepository)).
			ApproveBooking(ctx, testStudent(), 10, now)
		assert.ErrorIs(t, err, errors.ErrForbidden)
	})

	t.Run("success", func(t *testing.T) {
		bookings := new(mockBookingRepository)
		events := new(mockEventRepository)
		spaces := new(mockSpaceRepository)
		users := new(mockUserRepository)
		notify := new(mockNotifyRepository)

		bookings.On("GetBookingByID", ctx, int64(10)).Return(pending, nil)
		bookings.On("ListOverlappingBookings", ctx, pending.SpaceID, pending.StartsAt, pending.EndsAt, &exclude).Return([]entities.Booking{}, nil)
		events.On("ListOverlappingEvents", ctx, pending.SpaceID, pending.StartsAt, pending.EndsAt, (*int64)(nil)).Return([]entities.Event{}, nil)
		bookings.On("ApproveBooking", ctx, int64(10), testAdmin().UserID, now).Return(&approved, nil)
		spaces.On("GetSpaceByID", ctx, pending.SpaceID).Return(&hall, nil)
		users.On("GetUserByID", ctx, owner.UserID).Return(&owner, nil)
		notify.On("Notify", mock.Anything, []int64{owner.VKUserID}, mock.AnythingOfType("string")).Return([]int64{}, nil)

		got, err := NewBookingService(bookings, events, spaces, users, notify).ApproveBooking(ctx, testAdmin(), 10, now)
		assert.NoError(t, err)
		assert.Equal(t, &approved, got)
		bookings.AssertExpectations(t)
		notify.AssertExpectations(t)
	})
}
