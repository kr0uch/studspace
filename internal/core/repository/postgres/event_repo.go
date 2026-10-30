package postgres

import (
	"context"
	"time"

	"github.com/kr0uch/studspace/internal/core/models/dao"
	"github.com/kr0uch/studspace/internal/core/models/entities"
	dbpg "github.com/kr0uch/studspace/pkg/database/postgres"
	"github.com/kr0uch/studspace/pkg/errors"
)

type EventRepository struct {
	postgresDB *dbpg.PostgresDB
}

func NewEventRepository(postgresDB *dbpg.PostgresDB) *EventRepository {
	return &EventRepository{postgresDB: postgresDB}
}

func (r *EventRepository) GetEventView(ctx context.Context, eventID, userID int64) (*entities.EventView, error) {
	var row dao.EventView
	if err := r.postgresDB.GetContext(ctx, &row, qEventGetView, eventID, userID); err != nil {
		return nil, wrapErr(err)
	}
	view := row.ToEntity()
	return &view, nil
}

func (r *EventRepository) ListEventViews(ctx context.Context, from, to time.Time, userID int64, category *entities.EventCategory, mine bool) ([]entities.EventView, error) {
	categoryFilter := ""
	if category != nil {
		categoryFilter = string(*category)
	}

	var rows []dao.EventView
	if err := r.postgresDB.SelectContext(ctx, &rows, qEventListViews, from, to, userID, categoryFilter, mine); err != nil {
		return nil, wrapErr(err)
	}

	result := make([]entities.EventView, 0, len(rows))
	for _, row := range rows {
		result = append(result, row.ToEntity())
	}
	return result, nil
}

func (r *EventRepository) ListEventsInRange(ctx context.Context, from, to time.Time, spaceID *int16) ([]entities.Event, error) {
	var rows []dao.Event
	if err := r.postgresDB.SelectContext(ctx, &rows, qEventListInRange, from, to, spaceID); err != nil {
		return nil, wrapErr(err)
	}
	result := make([]entities.Event, 0, len(rows))
	for _, row := range rows {
		result = append(result, row.ToEntity())
	}
	return result, nil
}

func (r *EventRepository) ListBusyEventDays(ctx context.Context, from, to time.Time) ([]time.Time, error) {
	days := make([]time.Time, 0)
	if err := r.postgresDB.SelectContext(ctx, &days, qEventBusyDays, from, to); err != nil {
		return nil, wrapErr(err)
	}
	return days, nil
}

func (r *EventRepository) ListOverlappingEvents(ctx context.Context, spaceID int16, startsAt, endsAt time.Time, excludeEventID *int64) ([]entities.Event, error) {
	var rows []dao.Event
	if err := r.postgresDB.SelectContext(ctx, &rows, qEventListOverlapping, spaceID, startsAt, endsAt, excludeEventID); err != nil {
		return nil, wrapErr(err)
	}
	result := make([]entities.Event, 0, len(rows))
	for _, row := range rows {
		result = append(result, row.ToEntity())
	}
	return result, nil
}

func (r *EventRepository) CreateEvent(ctx context.Context, event entities.Event) (*entities.Event, error) {
	status := event.Status
	if status == "" {
		status = entities.EventStatusPublished
	}

	var row dao.Event
	err := r.postgresDB.GetContext(
		ctx,
		&row,
		qEventCreate,
		event.Title,
		event.Organizer,
		string(event.Category),
		event.SpaceID,
		event.StartsAt,
		event.EndsAt,
		event.Capacity,
		event.Description,
		string(status),
		event.CreatedBy,
	)
	if err != nil {
		return nil, wrapErr(err)
	}
	result := row.ToEntity()
	return &result, nil
}

func (r *EventRepository) UpdateEvent(ctx context.Context, event entities.Event) (*entities.Event, error) {
	var row dao.Event
	err := r.postgresDB.GetContext(
		ctx,
		&row,
		qEventUpdate,
		event.EventID,
		event.Title,
		event.Organizer,
		string(event.Category),
		event.SpaceID,
		event.StartsAt,
		event.EndsAt,
		event.Capacity,
		event.Description,
	)
	if err != nil {
		return nil, wrapErr(err)
	}
	result := row.ToEntity()
	return &result, nil
}

func (r *EventRepository) CancelEvent(ctx context.Context, eventID int64) (*entities.Event, error) {
	var row dao.Event
	if err := r.postgresDB.GetContext(ctx, &row, qEventCancel, eventID, string(entities.EventStatusCancelled)); err != nil {
		return nil, wrapErr(err)
	}
	result := row.ToEntity()
	return &result, nil
}

func (r *EventRepository) RegisterForEvent(ctx context.Context, eventID, userID int64) (*entities.EventView, error) {
	view, err := r.GetEventView(ctx, eventID, userID)
	if err != nil {
		return nil, err
	}
	if view.Status != entities.EventStatusPublished {
		return nil, errors.ErrNotFound
	}
	if view.IsRegistered {
		return view, nil
	}
	if view.SeatsLeft() == 0 {
		return nil, errors.ErrEventFull
	}

	if _, err := r.postgresDB.ExecContext(ctx, qEventUpsertRegistration, eventID, userID); err != nil {
		return nil, wrapErr(err)
	}
	if _, err := r.postgresDB.ExecContext(ctx, qEventRemoveWaitlist, eventID, userID); err != nil {
		return nil, wrapErr(err)
	}
	return r.GetEventView(ctx, eventID, userID)
}

func (r *EventRepository) UnregisterFromEvent(ctx context.Context, eventID, userID int64) (*entities.EventView, error) {
	if _, err := r.postgresDB.ExecContext(ctx, qEventCancelRegistration, eventID, userID); err != nil {
		return nil, wrapErr(err)
	}
	return r.GetEventView(ctx, eventID, userID)
}

func (r *EventRepository) AddEventWaitlist(ctx context.Context, eventID, userID int64) error {
	_, err := r.postgresDB.ExecContext(ctx, qEventAddWaitlist, eventID, userID)
	return wrapErr(err)
}

func (r *EventRepository) RemoveEventWaitlist(ctx context.Context, eventID, userID int64) error {
	_, err := r.postgresDB.ExecContext(ctx, qEventRemoveWaitlist, eventID, userID)
	return wrapErr(err)
}

func (r *EventRepository) ClaimEventWaitlist(ctx context.Context, eventID int64) ([]entities.EventWaitlistEntry, error) {
	var rows []dao.EventWaitlistEntry
	if err := r.postgresDB.SelectContext(ctx, &rows, qEventClaimWaitlist, eventID); err != nil {
		return nil, wrapErr(err)
	}
	result := make([]entities.EventWaitlistEntry, 0, len(rows))
	for _, row := range rows {
		result = append(result, row.ToEntity())
	}
	return result, nil
}
