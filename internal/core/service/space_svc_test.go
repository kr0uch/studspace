package service

import (
	"context"
	"testing"
	"time"

	"github.com/kr0uch/studspace/internal/core/models/dto"
	"github.com/kr0uch/studspace/internal/core/models/entities"
	"github.com/kr0uch/studspace/pkg/errors"
	"github.com/kr0uch/studspace/pkg/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSpaceService_ListActiveSpaces(t *testing.T) {
	ctx := context.Background()
	halls := []entities.Space{testHall()}

	tests := []struct {
		name          string
		mockBehavior  func(spaces *mockSpaceRepository)
		expectedError error
	}{
		{
			name: "success",
			mockBehavior: func(spaces *mockSpaceRepository) {
				spaces.On("ListActiveSpaces", ctx).Return(halls, nil)
			},
		},
		{
			name: "failure (db)",
			mockBehavior: func(spaces *mockSpaceRepository) {
				spaces.On("ListActiveSpaces", ctx).Return(nil, errors.ErrInternal)
			},
			expectedError: errors.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			spaces := new(mockSpaceRepository)
			tt.mockBehavior(spaces)

			got, err := NewSpaceService(spaces, new(mockBookingRepository), new(mockEventRepository)).ListActiveSpaces(ctx)
			if tt.expectedError != nil {
				assert.ErrorIs(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, halls, got)
			}
			spaces.AssertExpectations(t)
		})
	}
}

func TestSpaceService_GetSpaceByID(t *testing.T) {
	ctx := context.Background()
	hall := testHall()
	inactive := hall
	inactive.IsActive = false

	tests := []struct {
		name          string
		mockBehavior  func(spaces *mockSpaceRepository)
		expectedError error
	}{
		{
			name: "success",
			mockBehavior: func(spaces *mockSpaceRepository) {
				spaces.On("GetSpaceByID", ctx, int16(1)).Return(&hall, nil)
			},
		},
		{
			name: "failure (inactive)",
			mockBehavior: func(spaces *mockSpaceRepository) {
				spaces.On("GetSpaceByID", ctx, int16(1)).Return(&inactive, nil)
			},
			expectedError: errors.ErrNotFound,
		},
		{
			name: "failure (not found)",
			mockBehavior: func(spaces *mockSpaceRepository) {
				spaces.On("GetSpaceByID", ctx, int16(1)).Return(nil, errors.ErrNotFound)
			},
			expectedError: errors.ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			spaces := new(mockSpaceRepository)
			tt.mockBehavior(spaces)

			got, err := NewSpaceService(spaces, new(mockBookingRepository), new(mockEventRepository)).GetSpaceByID(ctx, 1)
			if tt.expectedError != nil {
				assert.ErrorIs(t, err, tt.expectedError)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, &hall, got)
			}
			spaces.AssertExpectations(t)
		})
	}
}

func TestSpaceService_GetSpaceCurrentStatus(t *testing.T) {
	ctx := context.Background()
	hall := testHall()
	now := utils.AtClock(time.Date(2026, 11, 2, 0, 0, 0, 0, utils.Moscow), 10*time.Hour)
	from, to := utils.DayRange(now)
	spaceID := hall.SpaceID

	spaces := new(mockSpaceRepository)
	spaces.On("GetSpaceByID", ctx, hall.SpaceID).Return(&hall, nil)

	bookings := new(mockBookingRepository)
	bookings.On("ListBookingsInRange", ctx, from, to, mock.MatchedBy(func(id *int16) bool {
		return id != nil && *id == spaceID
	})).Return([]entities.Booking{}, nil)

	events := new(mockEventRepository)
	events.On("ListEventsInRange", ctx, from, to, mock.MatchedBy(func(id *int16) bool {
		return id != nil && *id == spaceID
	})).Return([]entities.Event{}, nil)

	got, err := NewSpaceService(spaces, bookings, events).GetSpaceCurrentStatus(ctx, dto.GetSpaceCurrentStatus{
		SpaceID: hall.SpaceID,
		Now:     now,
	})
	assert.NoError(t, err)
	assert.Equal(t, entities.SpaceStateFree, got.State)
	spaces.AssertExpectations(t)
	bookings.AssertExpectations(t)
	events.AssertExpectations(t)
}

func TestSpaceService_GetSpaceAvailability_validation(t *testing.T) {
	_, err := NewSpaceService(new(mockSpaceRepository), new(mockBookingRepository), new(mockEventRepository)).
		GetSpaceAvailability(context.Background(), dto.GetSpaceAvailability{})
	assert.ErrorIs(t, err, errors.ErrValidation)
}
