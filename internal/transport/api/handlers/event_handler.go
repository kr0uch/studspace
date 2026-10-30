package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/kr0uch/studspace/internal/core/models/api"
	"github.com/kr0uch/studspace/internal/core/models/dto"
	"github.com/kr0uch/studspace/internal/core/models/entities"
	"github.com/kr0uch/studspace/pkg/utils"
	"github.com/kr0uch/studspace/pkg/web"
)

type EventService interface {
	ListEventViews(ctx context.Context, userID int64, in dto.GetEventViews) ([]entities.EventView, error)
	RegisterForEvent(ctx context.Context, userID, eventID int64) (*entities.EventView, error)
	UnregisterFromEvent(ctx context.Context, userID, eventID int64) (*entities.EventView, error)
	AddEventWaitlist(ctx context.Context, userID, eventID int64) error
	RemoveEventWaitlist(ctx context.Context, userID, eventID int64) error
}

type EventSpaceService interface {
	ListActiveSpaces(ctx context.Context) ([]entities.Space, error)
	GetSpaceByID(ctx context.Context, spaceID int16) (*entities.Space, error)
}

type EventController struct {
	eventService EventService
	spaceService EventSpaceService
}

func NewEventController(eventService EventService, spaceService EventSpaceService) *EventController {
	return &EventController{eventService: eventService, spaceService: spaceService}
}

func (c *EventController) ListEventsHandler(w http.ResponseWriter, r *http.Request) error {
	user, err := web.GetUserFromCtx(r.Context())
	if err != nil {
		return err
	}

	from, to, err := parseInclusiveDays(r.URL.Query().Get("from"), r.URL.Query().Get("to"))
	if err != nil {
		return err
	}

	var category *entities.EventCategory
	if raw := r.URL.Query().Get("category"); raw != "" {
		value := entities.EventCategory(raw)
		category = &value
	}

	now := time.Now()
	views, err := c.eventService.ListEventViews(r.Context(), user.UserID, dto.GetEventViews{
		From:     from,
		To:       to,
		Category: category,
		Mine:     r.URL.Query().Get("mine") == "true",
	})
	if err != nil {
		return err
	}

	spaces, err := c.spaceService.ListActiveSpaces(r.Context())
	if err != nil {
		return err
	}

	items := make([]api.Event, 0, len(views))
	for _, view := range views {
		items = append(items, api.NewEventView(view, spaceName(spaces, view.SpaceID), now))
	}

	return json.NewEncoder(w).Encode(api.EventList{
		WeekTotal: len(items),
		Featured:  featuredToday(items, now),
		Items:     items,
	})
}

func (c *EventController) RegisterForEventHandler(w http.ResponseWriter, r *http.Request) error {
	user, err := web.GetUserFromCtx(r.Context())
	if err != nil {
		return err
	}

	eventID, err := parseInt64ID(r.PathValue("id"), "id")
	if err != nil {
		return err
	}

	view, err := c.eventService.RegisterForEvent(r.Context(), user.UserID, eventID)
	if err != nil {
		return err
	}

	out, err := c.eventView(r, *view)
	if err != nil {
		return err
	}
	w.WriteHeader(http.StatusCreated)
	return json.NewEncoder(w).Encode(out)
}

func (c *EventController) UnregisterFromEventHandler(w http.ResponseWriter, r *http.Request) error {
	user, err := web.GetUserFromCtx(r.Context())
	if err != nil {
		return err
	}

	eventID, err := parseInt64ID(r.PathValue("id"), "id")
	if err != nil {
		return err
	}

	view, err := c.eventService.UnregisterFromEvent(r.Context(), user.UserID, eventID)
	if err != nil {
		return err
	}

	out, err := c.eventView(r, *view)
	if err != nil {
		return err
	}
	return json.NewEncoder(w).Encode(out)
}

func (c *EventController) AddWaitlistHandler(w http.ResponseWriter, r *http.Request) error {
	user, err := web.GetUserFromCtx(r.Context())
	if err != nil {
		return err
	}

	eventID, err := parseInt64ID(r.PathValue("id"), "id")
	if err != nil {
		return err
	}

	if err := c.eventService.AddEventWaitlist(r.Context(), user.UserID, eventID); err != nil {
		return err
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}

func (c *EventController) RemoveWaitlistHandler(w http.ResponseWriter, r *http.Request) error {
	user, err := web.GetUserFromCtx(r.Context())
	if err != nil {
		return err
	}

	eventID, err := parseInt64ID(r.PathValue("id"), "id")
	if err != nil {
		return err
	}

	if err := c.eventService.RemoveEventWaitlist(r.Context(), user.UserID, eventID); err != nil {
		return err
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}

func (c *EventController) eventView(r *http.Request, view entities.EventView) (api.Event, error) {
	name := ""
	if space, err := c.spaceService.GetSpaceByID(r.Context(), view.SpaceID); err == nil {
		name = space.Name
	}
	return api.NewEventView(view, name, time.Now()), nil
}

func featuredToday(items []api.Event, now time.Time) *api.Event {
	local := now.In(utils.Moscow)
	year, month, day := local.Date()
	for i := range items {
		if items[i].State == entities.EventUserStatePast {
			continue
		}
		starts := items[i].StartsAt.In(utils.Moscow)
		y, m, d := starts.Date()
		if y == year && m == month && d == day {
			return &items[i]
		}
	}
	return nil
}
