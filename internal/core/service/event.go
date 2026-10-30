package service

import (
	"context"
	"strings"
	"time"

	"github.com/kr0uch/studspace/internal/core/models/dto"
	"github.com/kr0uch/studspace/internal/core/models/entities"
	"github.com/kr0uch/studspace/internal/core/repository/postgres"
	"github.com/kr0uch/studspace/pkg/errors"
	"github.com/kr0uch/studspace/pkg/utils"
)

type EventService struct {
	events   *postgres.EventRepository
	bookings *postgres.BookingRepository
	spaces   *postgres.SpaceRepository
}

func NewEventService(
	events *postgres.EventRepository,
	bookings *postgres.BookingRepository,
	spaces *postgres.SpaceRepository,
) *EventService {
	return &EventService{events: events, bookings: bookings, spaces: spaces}
}

func (s *EventService) GetEventView(ctx context.Context, userID, eventID int64) (*entities.EventView, error) {
	return s.events.GetEventView(ctx, eventID, userID)
}

func (s *EventService) ListEventViews(ctx context.Context, userID int64, in dto.GetEventViews) ([]entities.EventView, error) {
	if in.Category != nil && !in.Category.Valid() {
		return nil, errors.Validation(map[string]string{"category": "неизвестная категория"})
	}
	return s.events.ListEventViews(ctx, in.From, in.To, userID, in.Category, in.Mine)
}

func (s *EventService) ListEventsInRange(ctx context.Context, from, to time.Time) ([]entities.Event, error) {
	if from.IsZero() || to.IsZero() || !to.After(from) {
		return nil, errors.Validation(map[string]string{"from": "укажите корректный интервал"})
	}
	return s.events.ListEventsInRange(ctx, from, to, nil)
}

func (s *EventService) ListBusyEventDays(ctx context.Context, month time.Time) ([]time.Time, error) {
	if month.IsZero() {
		return nil, errors.Validation(map[string]string{"month": "обязательное поле"})
	}
	from, to := utils.MonthRange(month)
	return s.events.ListBusyEventDays(ctx, from, to)
}

func (s *EventService) RegisterForEvent(ctx context.Context, userID, eventID int64) (*entities.EventView, error) {
	return s.events.RegisterForEvent(ctx, eventID, userID)
}

func (s *EventService) UnregisterFromEvent(ctx context.Context, userID, eventID int64) (*entities.EventView, error) {
	if _, err := s.events.UnregisterFromEvent(ctx, eventID, userID); err != nil {
		return nil, err
	}
	if _, err := s.events.ClaimEventWaitlist(ctx, eventID); err != nil {
		return nil, err
	}
	return s.events.GetEventView(ctx, eventID, userID)
}

func (s *EventService) AddEventWaitlist(ctx context.Context, userID, eventID int64) error {
	view, err := s.events.GetEventView(ctx, eventID, userID)
	if err != nil {
		return err
	}
	if view.Status != entities.EventStatusPublished {
		return errors.ErrNotFound
	}
	if view.IsRegistered {
		return errors.ErrBadRequest.WithMessage("вы уже записаны")
	}
	if view.SeatsLeft() > 0 {
		return errors.ErrBadRequest.WithMessage("есть свободные места")
	}
	return s.events.AddEventWaitlist(ctx, eventID, userID)
}

func (s *EventService) RemoveEventWaitlist(ctx context.Context, userID, eventID int64) error {
	return s.events.RemoveEventWaitlist(ctx, eventID, userID)
}

func (s *EventService) CreateEvent(ctx context.Context, actor entities.User, in dto.CreateEvent, now time.Time) (*entities.Event, error) {
	if err := requireAdmin(actor); err != nil {
		return nil, err
	}
	in, err := s.prepareEvent(ctx, in, now)
	if err != nil {
		return nil, err
	}

	overlaps, err := loadSlotOverlaps(ctx, s.bookings, s.events, in.SpaceID, in.StartsAt, in.EndsAt, nil, nil)
	if err != nil {
		return nil, err
	}
	if overlaps.HardConflict() {
		return nil, errors.ErrSlotConflict
	}

	return s.events.CreateEvent(ctx, entities.Event{
		Title:       in.Title,
		Organizer:   in.Organizer,
		Category:    in.Category,
		SpaceID:     in.SpaceID,
		StartsAt:    in.StartsAt,
		EndsAt:      in.EndsAt,
		Capacity:    in.Capacity,
		Description: in.Description,
		Status:      entities.EventStatusPublished,
		CreatedBy:   actor.UserID,
	})
}

func (s *EventService) UpdateEvent(ctx context.Context, actor entities.User, eventID int64, in dto.CreateEvent, now time.Time) (*entities.Event, error) {
	if err := requireAdmin(actor); err != nil {
		return nil, err
	}
	current, err := s.events.GetEventView(ctx, eventID, actor.UserID)
	if err != nil {
		return nil, err
	}
	if current.Status != entities.EventStatusPublished {
		return nil, errors.ErrInvalidTransition
	}

	in, err = s.prepareEvent(ctx, in, now)
	if err != nil {
		return nil, err
	}
	if in.Capacity < current.RegisteredCount {
		return nil, errors.Validation(map[string]string{"capacity": "не меньше числа записавшихся"})
	}

	exclude := eventID
	overlaps, err := loadSlotOverlaps(ctx, s.bookings, s.events, in.SpaceID, in.StartsAt, in.EndsAt, nil, &exclude)
	if err != nil {
		return nil, err
	}
	if overlaps.HardConflict() {
		return nil, errors.ErrSlotConflict
	}

	return s.events.UpdateEvent(ctx, entities.Event{
		EventID:     eventID,
		Title:       in.Title,
		Organizer:   in.Organizer,
		Category:    in.Category,
		SpaceID:     in.SpaceID,
		StartsAt:    in.StartsAt,
		EndsAt:      in.EndsAt,
		Capacity:    in.Capacity,
		Description: in.Description,
	})
}

func (s *EventService) prepareEvent(ctx context.Context, in dto.CreateEvent, now time.Time) (dto.CreateEvent, error) {
	space, err := activeSpace(ctx, s.spaces, in.SpaceID)
	if err != nil {
		return dto.CreateEvent{}, err
	}
	if err := validateSpaceTimeWindow(*space, in.StartsAt, in.EndsAt, now); err != nil {
		return dto.CreateEvent{}, err
	}
	in.Title = strings.TrimSpace(in.Title)
	in.Organizer = strings.TrimSpace(in.Organizer)

	fields := map[string]string{}
	if in.Title == "" {
		fields["title"] = "обязательное поле"
	}
	if in.Organizer == "" {
		fields["organizer"] = "обязательное поле"
	}
	if !in.Category.Valid() {
		fields["category"] = "неизвестная категория"
	}
	if in.Capacity < 1 || in.Capacity > space.Capacity {
		fields["capacity"] = "не больше вместимости помещения"
	}
	if len(fields) > 0 {
		return dto.CreateEvent{}, errors.Validation(fields)
	}
	return in, nil
}

func (s *EventService) CancelEvent(ctx context.Context, actor entities.User, eventID int64) (*entities.Event, error) {
	if err := requireAdmin(actor); err != nil {
		return nil, err
	}
	return s.events.CancelEvent(ctx, eventID)
}
