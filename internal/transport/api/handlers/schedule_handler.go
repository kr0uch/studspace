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

type ScheduleSpaceService interface {
	ListActiveSpaces(ctx context.Context) ([]entities.Space, error)
}

type ScheduleController struct {
	spaceService   ScheduleSpaceService
	bookingService ScheduleBookingService
}

func NewScheduleController(spaceService ScheduleSpaceService, bookingService ScheduleBookingService) *ScheduleController {
	return &ScheduleController{
		spaceService:   spaceService,
		bookingService: bookingService,
	}
}

// GetScheduleHandler godoc
// @Summary Расписание
// @Description Брони в диапазоне дат: свои и чужие pending и approved.
// @Tags Расписание
// @Produce json
// @Security VKAuth
// @Param from query string true "YYYY-MM-DD"
// @Param to query string true "YYYY-MM-DD"
// @Success 200 {object} api.Schedule
// @Failure 401 {object} errors.ErrorResponse
// @Failure 422 {object} errors.ErrorResponse
// @Failure 500 {object} errors.ErrorResponse
// @Router /schedule [get]
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

	now := time.Now()
	items := make([]api.ScheduleItem, 0, len(bookings))
	for _, booking := range bookings {
		isMine := booking.UserID == user.UserID
		items = append(items, api.ScheduleItem{
			BookingID: booking.BookingID,
			SpaceID:   booking.SpaceID,
			Title:     booking.Title,
			StartsAt:  booking.StartsAt,
			EndsAt:    booking.EndsAt,
			OwnerName: booking.FirstName + " " + booking.LastName,
			Status:    booking.CurrentStatus(now),
			IsMine:    isMine,
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

// GetBusyDaysHandler godoc
// @Summary Занятые дни календаря
// @Description Дни месяца с pending или approved бронями.
// @Tags Расписание
// @Produce json
// @Security VKAuth
// @Param month query string true "YYYY-MM"
// @Success 200 {object} api.BusyDays
// @Failure 401 {object} errors.ErrorResponse
// @Failure 422 {object} errors.ErrorResponse
// @Failure 500 {object} errors.ErrorResponse
// @Router /calendar/busy-days [get]
func (c *ScheduleController) GetBusyDaysHandler(w http.ResponseWriter, r *http.Request) error {
	month, err := parseMonth(r.URL.Query().Get("month"))
	if err != nil {
		return err
	}

	days, err := c.bookingService.ListBusyBookingDays(r.Context(), month)
	if err != nil {
		return err
	}

	return json.NewEncoder(w).Encode(api.BusyDays{
		Month: month,
		Days:  uniqueDays(days),
	})
}

func uniqueDays(part []time.Time) []time.Time {
	seen := make(map[time.Time]bool, len(part))
	days := make([]time.Time, 0, len(part))
	for _, day := range part {
		if seen[day] {
			continue
		}
		seen[day] = true
		days = append(days, day)
	}
	sort.Slice(days, func(i, j int) bool {
		return days[i].Before(days[j])
	})
	return days
}
