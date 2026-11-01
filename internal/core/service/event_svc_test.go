package service

import (
	"context"
	"testing"
	"time"

	"github.com/kr0uch/studspace/internal/core/models/dto"
	"github.com/kr0uch/studspace/internal/core/models/entities"
	"github.com/kr0uch/studspace/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestEventService_AddEventWaitlist(t *testing.T) {
	ctx := context.Background()
	full := &entities.EventView{
		Event:           entities.Event{EventID: 5, Status: entities.EventStatusPublished, Capacity: 1},
		RegisteredCount: 1,
	}
	open := &entities.EventView{
		Event:           entities.Event{EventID: 5, Status: entities.EventStatusPublished, Capacity: 10},
		RegisteredCount: 1,
	}
	cancelled := &entities.EventView{
		Event: entities.Event{EventID: 5, Status: entities.EventStatusCancelled, Capacity: 10},
	}

	tests := []struct {
		name          string
		mockBehavior  func(events *mockEventRepository)
		expectedError error
	}{
		{
			name: "success",
			mockBehavior: func(events *mockEventRepository) {
				events.On("GetEventView", ctx, int64(5), int64(1)).Return(full, nil)
				events.On("AddEventWaitlist", ctx, int64(5), int64(1)).Return(nil)
			},
		},
		{
			name: "failure (seats left)",
			mockBehavior: func(events *mockEventRepository) {
				events.On("GetEventView", ctx, int64(5), int64(1)).Return(open, nil)
			},
			expectedError: errors.ErrBadRequest,
		},
		{
			name: "failure (cancelled)",
			mockBehavior: func(events *mockEventRepository) {
				events.On("GetEventView", ctx, int64(5), int64(1)).Return(cancelled, nil)
			},
			expectedError: errors.ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			events := new(mockEventRepository)
			tt.mockBehavior(events)

			err := NewEventService(events, new(mockBookingRepository), new(mockSpaceRepository), new(mockUserRepository), new(mockNotifyRepository)).
				AddEventWaitlist(ctx, 1, 5)
			if tt.expectedError != nil {
				assert.ErrorIs(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err)
			}
			events.AssertExpectations(t)
		})
	}
}

func TestEventService_RegisterForEvent(t *testing.T) {
	ctx := context.Background()
	view := &entities.EventView{Event: entities.Event{EventID: 5, Status: entities.EventStatusPublished, Capacity: 10}}
	events := new(mockEventRepository)
	events.On("RegisterForEvent", ctx, int64(5), int64(1)).Return(view, nil)

	got, err := NewEventService(events, new(mockBookingRepository), new(mockSpaceRepository), new(mockUserRepository), new(mockNotifyRepository)).
		RegisterForEvent(ctx, 1, 5)
	assert.NoError(t, err)
	assert.Equal(t, view, got)
	events.AssertExpectations(t)
}

func TestEventService_CreateEvent_forbidden(t *testing.T) {
	_, err := NewEventService(new(mockEventRepository), new(mockBookingRepository), new(mockSpaceRepository), new(mockUserRepository), new(mockNotifyRepository)).
		CreateEvent(context.Background(), testStudent(), dto.CreateEvent{}, time.Now())
	assert.ErrorIs(t, err, errors.ErrForbidden)
}
