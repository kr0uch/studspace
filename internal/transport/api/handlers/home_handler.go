package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"sort"
	"time"

	"github.com/kr0uch/studspace/internal/core/models/api"
	"github.com/kr0uch/studspace/internal/core/models/dto"
	"github.com/kr0uch/studspace/internal/core/models/entities"
	"github.com/kr0uch/studspace/pkg/utils"
	"github.com/kr0uch/studspace/pkg/web"
)

type HomeBookingService interface {
	ListUserBookings(ctx context.Context, actor entities.User, in dto.ListUserBookings) ([]entities.Booking, *string, error)
}

type HomeEventService interface {
	ListEventViews(ctx context.Context, userID int64, in dto.GetEventViews) ([]entities.EventView, error)
}

type HomeSpaceService interface {
	ListActiveSpaces(ctx context.Context) ([]entities.Space, error)
	GetSpaceCurrentStatus(ctx context.Context, in dto.GetSpaceCurrentStatus) (*entities.SpaceNowStatus, error)
}

type HomeController struct {
	bookingService HomeBookingService
	eventService   HomeEventService
	spaceService   HomeSpaceService
}

func NewHomeController(bookingService HomeBookingService, eventService HomeEventService, spaceService HomeSpaceService) *HomeController {
	return &HomeController{
		bookingService: bookingService,
		eventService:   eventService,
		spaceService:   spaceService,
	}
}

// GetHomeHandler godoc
// @Summary Главная
// @Description Сегодняшние активные брони и свои мероприятия, плюс статусы залов.
// @Tags Главная
// @Produce json
// @Security VKAuth
// @Success 200 {object} api.Home
// @Failure 401 {object} errors.ErrorResponse
// @Failure 404 {object} errors.ErrorResponse
// @Failure 500 {object} errors.ErrorResponse
// @Router /home [get]
func (c *HomeController) GetHomeHandler(w http.ResponseWriter, r *http.Request) error {
	user, err := web.GetUserFromCtx(r.Context())
	if err != nil {
		return err
	}

	now := time.Now()
	from, to := utils.DayRange(now)

	bookings, _, err := c.bookingService.ListUserBookings(r.Context(), *user, dto.ListUserBookings{
		UserID: user.UserID,
		Tab:    entities.UserBookingTabActive,
		Now:    now,
	})
	if err != nil {
		return err
	}

	events, err := c.eventService.ListEventViews(r.Context(), user.UserID, dto.GetEventViews{
		From: from,
		To:   to,
		Mine: true,
	})
	if err != nil {
		return err
	}

	spaces, err := c.spaceService.ListActiveSpaces(r.Context())
	if err != nil {
		return err
	}

	spaceViews := make([]api.Space, 0, len(spaces))
	for _, space := range spaces {
		status, err := c.spaceService.GetSpaceCurrentStatus(r.Context(), dto.GetSpaceCurrentStatus{
			SpaceID: space.SpaceID,
			Now:     now,
		})
		if err != nil {
			return err
		}
		spaceViews = append(spaceViews, api.Space{Space: space, Now: status})
	}

	today := make([]api.HomeTodayItem, 0)
	for _, booking := range bookings {
		if !booking.StartsAt.Before(to) || !booking.EndsAt.After(from) {
			continue
		}
		today = append(today, api.HomeTodayItem{
			Type:      api.HomeItemBooking,
			BookingID: booking.BookingID,
			StartsAt:  booking.StartsAt,
			EndsAt:    booking.EndsAt,
			Space:     api.HomeTodaySpace{SpaceID: booking.SpaceID, Name: spaceName(spaces, booking.SpaceID)},
			Purpose:   booking.Purpose,
			Status:    booking.CurrentStatus(now),
		})
	}
	for _, event := range events {
		if !event.IsRegistered {
			continue
		}
		today = append(today, api.HomeTodayItem{
			Type:     api.HomeItemEvent,
			EventID:  event.EventID,
			StartsAt: event.StartsAt,
			EndsAt:   event.EndsAt,
			Space:    api.HomeTodaySpace{SpaceID: event.SpaceID, Name: spaceName(spaces, event.SpaceID)},
			Title:    event.Title,
		})
	}
	sort.Slice(today, func(i, j int) bool {
		if !today[i].StartsAt.Equal(today[j].StartsAt) {
			return today[i].StartsAt.Before(today[j].StartsAt)
		}
		if today[i].Type != today[j].Type {
			return today[i].Type < today[j].Type
		}
		return today[i].BookingID+today[i].EventID < today[j].BookingID+today[j].EventID
	})

	return json.NewEncoder(w).Encode(api.Home{
		Today:  api.HomeToday{Total: len(today), Items: today},
		Spaces: spaceViews,
	})
}
