package handlers

import (
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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockSpaceService struct{ mock.Mock }

func (m *mockSpaceService) ListActiveSpaces(ctx context.Context) ([]entities.Space, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entities.Space), args.Error(1)
}

func (m *mockSpaceService) GetSpaceAvailability(ctx context.Context, in dto.GetSpaceAvailability) (*dto.SpaceAvailability, error) {
	args := m.Called(ctx, in)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.SpaceAvailability), args.Error(1)
}

func TestSpaceController_ListSpacesHandler(t *testing.T) {
	spaces := []entities.Space{{SpaceID: 1, Name: "Переговорная", IsActive: true}}

	tests := []struct {
		name          string
		mockBehavior  func(svc *mockSpaceService)
		expectedError error
	}{
		{
			name: "success",
			mockBehavior: func(svc *mockSpaceService) {
				svc.On("ListActiveSpaces", mock.Anything).Return(spaces, nil)
			},
		},
		{
			name: "failure (db)",
			mockBehavior: func(svc *mockSpaceService) {
				svc.On("ListActiveSpaces", mock.Anything).Return(nil, errors.ErrInternal)
			},
			expectedError: errors.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := new(mockSpaceService)
			tt.mockBehavior(svc)
			w := httptest.NewRecorder()
			err := NewSpaceController(svc).ListSpacesHandler(w, httptest.NewRequest(http.MethodGet, "/spaces", nil))
			if tt.expectedError != nil {
				assert.ErrorIs(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err)
				var got []api.Space
				require.NoError(t, json.NewDecoder(w.Body).Decode(&got))
				assert.Equal(t, "Переговорная", got[0].Name)
			}
			svc.AssertExpectations(t)
		})
	}
}

func TestSpaceController_GetAvailabilityHandler(t *testing.T) {
	t.Run("failure (bad id)", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/spaces/x/availability?date=2026-11-01", nil)
		req.SetPathValue("id", "x")
		err := NewSpaceController(new(mockSpaceService)).GetAvailabilityHandler(httptest.NewRecorder(), req)
		assert.ErrorIs(t, err, errors.ErrValidation)
	})

	t.Run("success", func(t *testing.T) {
		date := time.Date(2026, 11, 1, 0, 0, 0, 0, time.FixedZone("MSK", 3*3600))
		svc := new(mockSpaceService)
		svc.On("GetSpaceAvailability", mock.Anything, mock.AnythingOfType("dto.GetSpaceAvailability")).
			Return(&dto.SpaceAvailability{Space: entities.Space{SpaceID: 1, SlotStepMin: 30}, Date: date}, nil)

		req := httptest.NewRequest(http.MethodGet, "/spaces/1/availability?date=2026-11-01", nil)
		req.SetPathValue("id", "1")
		w := httptest.NewRecorder()
		err := NewSpaceController(svc).GetAvailabilityHandler(w, req)
		assert.NoError(t, err)
		svc.AssertExpectations(t)
	})
}
