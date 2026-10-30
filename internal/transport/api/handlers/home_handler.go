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

type HomeBookingService interface {
	ListUserBookings(ctx context.Context, actor entities.User, in dto.ListUserBookings) ([]entities.Booking, *string, error)
}

type HomeSpaceService interface {
	ListActiveSpaces(ctx context.Context) ([]entities.Space, error)
	GetSpaceCurrentStatus(ctx context.Context, in dto.GetSpaceCurrentStatus) (*entities.SpaceNowStatus, error)
}

type HomeController struct {
	bookingService HomeBookingService
	spaceService   HomeSpaceService
}

func NewHomeController(bookingService HomeBookingService, spaceService HomeSpaceService) *HomeController {
	return &HomeController{
		bookingService: bookingService,
		spaceService:   spaceService,
	}
}

func (c *HomeController) GetHomeHandler(w http.ResponseWriter, r *http.Request) error {
	user, err := web.GetUserFromCtx(r.Context())
	if err != nil {
		return err
	}

	now := time.Now()
	bookings, _, err := c.bookingService.ListUserBookings(r.Context(), *user, dto.ListUserBookings{
		UserID: user.UserID,
		Tab:    entities.UserBookingTabActive,
		Now:    now,
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

	from, to := utils.DayRange(now)
	today := make([]api.HomeBooking, 0)
	for _, booking := range bookings {
		if !booking.StartsAt.Before(to) || !booking.EndsAt.After(from) {
			continue
		}
		today = append(today, api.HomeBooking{
			BookingID: booking.BookingID,
			StartsAt:  booking.StartsAt,
			EndsAt:    booking.EndsAt,
			Space:     api.HomeBookingSpace{SpaceID: booking.SpaceID, Name: spaceName(spaces, booking.SpaceID)},
			Purpose:   booking.Purpose,
			Status:    booking.CurrentStatus(now),
		})
	}

	return json.NewEncoder(w).Encode(api.Home{
		TodayBookings: api.HomeTodayBookings{Total: len(today), Items: today},
		Spaces:        spaceViews,
	})
}
