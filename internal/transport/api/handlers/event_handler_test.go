package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kr0uch/studspace/internal/core/models/dto"
	"github.com/kr0uch/studspace/internal/core/models/entities"
	"github.com/kr0uch/studspace/pkg/errors"
	"github.com/kr0uch/studspace/pkg/web"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockEventService struct{ mock.Mock }

func (m *mockEventService) ListEventViews(ctx context.Context, userID int64, in dto.GetEventViews) ([]entities.EventView, error) {
	args := m.Called(ctx, userID, in)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entities.EventView), args.Error(1)
}

func (m *mockEventService) RegisterForEvent(ctx context.Context, userID, eventID int64) (*entities.EventView, error) {
	args := m.Called(ctx, userID, eventID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.EventView), args.Error(1)
}

func (m *mockEventService) UnregisterFromEvent(ctx context.Context, userID, eventID int64) (*entities.EventView, error) {
	args := m.Called(ctx, userID, eventID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.EventView), args.Error(1)
}

func (m *mockEventService) AddEventWaitlist(ctx context.Context, userID, eventID int64) error {
	return m.Called(ctx, userID, eventID).Error(0)
}

func (m *mockEventService) RemoveEventWaitlist(ctx context.Context, userID, eventID int64) error {
	return m.Called(ctx, userID, eventID).Error(0)
}

type mockEventSpaceService struct{ mock.Mock }

func (m *mockEventSpaceService) ListActiveSpaces(ctx context.Context) ([]entities.Space, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entities.Space), args.Error(1)
}

func (m *mockEventSpaceService) GetSpaceByID(ctx context.Context, spaceID int16) (*entities.Space, error) {
	args := m.Called(ctx, spaceID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Space), args.Error(1)
}

func TestEventController_RegisterForEventHandler(t *testing.T) {
	user := testHandlerUser()
	view := &entities.EventView{Event: entities.Event{EventID: 5, SpaceID: 1, Title: "Митап", Capacity: 10}}

	t.Run("success", func(t *testing.T) {
		events := new(mockEventService)
		spaces := new(mockEventSpaceService)
		events.On("RegisterForEvent", mock.Anything, user.UserID, int64(5)).Return(view, nil)
		spaces.On("GetSpaceByID", mock.Anything, int16(1)).Return(&entities.Space{SpaceID: 1, Name: "A"}, nil)

		req := httptest.NewRequest(http.MethodPost, "/events/5/registration", nil)
		req.SetPathValue("id", "5")
		req = req.WithContext(web.ContextWithUser(req.Context(), user))
		w := httptest.NewRecorder()
		err := NewEventController(events, spaces).RegisterForEventHandler(w, req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, w.Code)
		events.AssertExpectations(t)
	})

	t.Run("failure (full)", func(t *testing.T) {
		events := new(mockEventService)
		events.On("RegisterForEvent", mock.Anything, user.UserID, int64(5)).Return(nil, errors.ErrEventFull)
		req := httptest.NewRequest(http.MethodPost, "/events/5/registration", nil)
		req.SetPathValue("id", "5")
		req = req.WithContext(web.ContextWithUser(req.Context(), user))
		err := NewEventController(events, new(mockEventSpaceService)).RegisterForEventHandler(httptest.NewRecorder(), req)
		assert.ErrorIs(t, err, errors.ErrEventFull)
	})
}

func TestEventController_AddWaitlistHandler(t *testing.T) {
	user := testHandlerUser()
	events := new(mockEventService)
	events.On("AddEventWaitlist", mock.Anything, user.UserID, int64(5)).Return(nil)
	req := httptest.NewRequest(http.MethodPost, "/events/5/waitlist", nil)
	req.SetPathValue("id", "5")
	req = req.WithContext(web.ContextWithUser(req.Context(), user))
	w := httptest.NewRecorder()
	err := NewEventController(events, new(mockEventSpaceService)).AddWaitlistHandler(w, req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, w.Code)
}
