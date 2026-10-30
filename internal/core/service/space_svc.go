package service

import (
	"context"
	"sort"
	"time"

	"github.com/guregu/null/v6"

	"github.com/kr0uch/studspace/internal/core/models/dto"
	"github.com/kr0uch/studspace/internal/core/models/entities"
	"github.com/kr0uch/studspace/internal/core/repository/postgres"
	"github.com/kr0uch/studspace/pkg/errors"
	"github.com/kr0uch/studspace/pkg/utils"
)

type SpaceService struct {
	spaces   *postgres.SpaceRepository
	bookings *postgres.BookingRepository
	events   *postgres.EventRepository
}

func NewSpaceService(
	spaces *postgres.SpaceRepository,
	bookings *postgres.BookingRepository,
	events *postgres.EventRepository,
) *SpaceService {
	return &SpaceService{spaces: spaces, bookings: bookings, events: events}
}

func (s *SpaceService) ListActiveSpaces(ctx context.Context) ([]entities.Space, error) {
	return s.spaces.ListActiveSpaces(ctx)
}

func (s *SpaceService) GetSpaceByID(ctx context.Context, spaceID int16) (*entities.Space, error) {
	return activeSpace(ctx, s.spaces, spaceID)
}

func (s *SpaceService) GetSpaceAvailability(ctx context.Context, in dto.GetSpaceAvailability) (*dto.SpaceAvailability, error) {
	if in.Date.IsZero() {
		return nil, errors.Validation(map[string]string{"date": "обязательное поле"})
	}

	space, err := activeSpace(ctx, s.spaces, in.SpaceID)
	if err != nil {
		return nil, err
	}

	from, to := utils.DayRange(in.Date.In(utils.Moscow))
	busy, err := s.busyInRange(ctx, from, to, &space.SpaceID)
	if err != nil {
		return nil, err
	}

	return &dto.SpaceAvailability{Space: *space, Date: from, Busy: busy}, nil
}

func (s *SpaceService) GetSpaceCurrentStatus(ctx context.Context, in dto.GetSpaceCurrentStatus) (*entities.SpaceNowStatus, error) {
	if in.Now.IsZero() {
		return nil, errors.Validation(map[string]string{"now": "обязательное поле"})
	}

	space, err := activeSpace(ctx, s.spaces, in.SpaceID)
	if err != nil {
		return nil, err
	}

	now := in.Now
	day := now.In(utils.Moscow)
	from, to := utils.DayRange(day)
	closeAt := utils.AtClock(day, space.CloseTime)
	openAt := utils.AtClock(day, space.OpenTime)
	if now.Before(openAt) {
		return &entities.SpaceNowStatus{State: entities.SpaceStateClosed, Until: null.TimeFrom(openAt)}, nil
	}
	if !now.Before(closeAt) {
		nextOpen := utils.AtClock(day.AddDate(0, 0, 1), space.OpenTime)
		return &entities.SpaceNowStatus{State: entities.SpaceStateClosed, Until: null.TimeFrom(nextOpen)}, nil
	}

	busy, err := s.busyInRange(ctx, from, to, &space.SpaceID)
	if err != nil {
		return nil, err
	}
	merged := mergeBusy(busy)

	for _, interval := range merged {
		if !now.Before(interval.StartsAt) && now.Before(interval.EndsAt) {
			return &entities.SpaceNowStatus{State: entities.SpaceStateBusy, Until: null.TimeFrom(interval.EndsAt)}, nil
		}
		if interval.StartsAt.After(now) {
			return &entities.SpaceNowStatus{State: entities.SpaceStateFree, Until: null.TimeFrom(interval.StartsAt)}, nil
		}
	}
	return &entities.SpaceNowStatus{State: entities.SpaceStateFree, Until: null.TimeFrom(closeAt)}, nil
}

func activeSpace(ctx context.Context, spaces *postgres.SpaceRepository, spaceID int16) (*entities.Space, error) {
	space, err := spaces.GetSpaceByID(ctx, spaceID)
	if err != nil {
		return nil, err
	}
	if !space.IsActive {
		return nil, errors.ErrNotFound
	}
	return space, nil
}

func (s *SpaceService) busyInRange(ctx context.Context, from, to time.Time, spaceID *int16) ([]dto.BusyInterval, error) {
	bookings, err := s.bookings.ListBookingsInRange(ctx, from, to, spaceID)
	if err != nil {
		return nil, err
	}
	events, err := s.events.ListEventsInRange(ctx, from, to, spaceID)
	if err != nil {
		return nil, err
	}

	result := make([]dto.BusyInterval, 0, len(bookings)+len(events))
	for _, booking := range bookings {
		if booking.Status != entities.BookingStatusApproved {
			continue
		}
		result = append(result, dto.BusyInterval{
			StartsAt: booking.StartsAt,
			EndsAt:   booking.EndsAt,
			Kind:     dto.BusyKindBooking,
		})
	}
	for _, event := range events {
		result = append(result, dto.BusyInterval{
			StartsAt: event.StartsAt,
			EndsAt:   event.EndsAt,
			Kind:     dto.BusyKindEvent,
		})
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].StartsAt.Before(result[j].StartsAt)
	})
	return result, nil
}

func mergeBusy(intervals []dto.BusyInterval) []dto.BusyInterval {
	if len(intervals) == 0 {
		return nil
	}
	sort.Slice(intervals, func(i, j int) bool {
		return intervals[i].StartsAt.Before(intervals[j].StartsAt)
	})
	merged := []dto.BusyInterval{intervals[0]}
	for _, interval := range intervals[1:] {
		last := &merged[len(merged)-1]
		if interval.StartsAt.After(last.EndsAt) {
			merged = append(merged, interval)
			continue
		}
		if interval.EndsAt.After(last.EndsAt) {
			last.EndsAt = interval.EndsAt
		}
	}
	return merged
}
