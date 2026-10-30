package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"sort"
	"time"

	"github.com/kr0uch/studspace/internal/core/models/api"
	"github.com/kr0uch/studspace/internal/core/models/entities"
	"github.com/kr0uch/studspace/pkg/web"
)

type ScheduleBookingService interface {
	ListBookingsInRange(ctx context.Context, from, to time.Time) ([]entities.Booking, error)
	ListBusyBookingDays(ctx context.Context, month time.Time) ([]time.Time, error)
}

type ScheduleEventService interface {
	ListEventsInRange(ctx context.Context, from, to time.Time) ([]entities.Event, error)
	ListBusyEventDays(ctx context.Context, month time.Time) ([]time.Time, error)
}

type ScheduleSpaceService interface {
	ListActiveSpaces(ctx context.Context) ([]entities.Space, error)
}

type ScheduleController struct {
	spaceService   ScheduleSpaceService
	bookingService ScheduleBookingService
	eventService   ScheduleEventService
}

func NewScheduleController(spaceService ScheduleSpaceService, bookingService ScheduleBookingService, eventService ScheduleEventService) *ScheduleController {
	return &ScheduleController{
		spaceService:   spaceService,
		bookingService: bookingService,
		eventService:   eventService,
	}
}

func (c *ScheduleController) GetScheduleHandler(w http.ResponseWriter, r *http.Request) error {
	user, err := web.GetUserFromCtx(r.Context())
	if err != nil {
		return err
	}

	from, to, err := parseInclusiveDays(r.URL.Query().Get("from"), r.URL.Query().Get("to"))
	if err != nil {
		return err
	}

	spaces, err := c.spaceService.ListActiveSpaces(r.Context())
	if err != nil {
		return err
	}

	bookings, err := c.bookingService.ListBookingsInRange(r.Context(), from, to)
	if err != nil {
		return err
	}

	events, err := c.eventService.ListEventsInRange(r.Context(), from, to)
	if err != nil {
		return err
	}

	now := time.Now()
	items := make([]api.ScheduleItem, 0, len(bookings)+len(events))
	for _, booking := range bookings {
		isMine := booking.UserID == user.UserID
		if booking.Status == entities.BookingStatusPending && !isMine {
			continue
		}
		items = append(items, api.ScheduleItem{
			Type:      "booking",
			BookingID: booking.BookingID,
			SpaceID:   booking.SpaceID,
			Title:     booking.Purpose,
			StartsAt:  booking.StartsAt,
			EndsAt:    booking.EndsAt,
			OwnerName: booking.FirstName + " " + booking.LastName,
			Status:    booking.CurrentStatus(now),
			IsMine:    isMine,
		})
	}
	for _, event := range events {
		items = append(items, api.ScheduleItem{
			Type:      "event",
			EventID:   event.EventID,
			SpaceID:   event.SpaceID,
			Title:     event.Title,
			StartsAt:  event.StartsAt,
			EndsAt:    event.EndsAt,
			Organizer: event.Organizer,
		})
	}

	spaceViews := make([]api.ScheduleSpace, 0, len(spaces))
	for _, space := range spaces {
		spaceViews = append(spaceViews, api.ScheduleSpace{
			SpaceID:  space.SpaceID,
			Name:     space.Name,
			Capacity: space.Capacity,
			Color:    space.Color,
		})
	}

	return json.NewEncoder(w).Encode(api.Schedule{Spaces: spaceViews, Items: items})
}

func (c *ScheduleController) GetBusyDaysHandler(w http.ResponseWriter, r *http.Request) error {
	month, err := parseMonth(r.URL.Query().Get("month"))
	if err != nil {
		return err
	}

	bookingDays, err := c.bookingService.ListBusyBookingDays(r.Context(), month)
	if err != nil {
		return err
	}
	eventDays, err := c.eventService.ListBusyEventDays(r.Context(), month)
	if err != nil {
		return err
	}

	return json.NewEncoder(w).Encode(api.BusyDays{
		Month: month,
		Days:  mergeDays(bookingDays, eventDays),
	})
}

func mergeDays(parts ...[]time.Time) []time.Time {
	seen := make(map[time.Time]bool)
	for _, part := range parts {
		for _, day := range part {
			if !seen[day] {
				seen[day] = true
			}
		}
	}
	days := make([]time.Time, 0, len(seen))
	for day := range seen {
		days = append(days, day)
	}
	sort.Slice(days, func(i, j int) bool {
		return days[i].Before(days[j])
	})
	return days
}
