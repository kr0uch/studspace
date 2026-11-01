package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/kr0uch/studspace/internal/core/models/api"
	"github.com/kr0uch/studspace/internal/core/models/dto"
	"github.com/kr0uch/studspace/internal/core/models/entities"
	"github.com/kr0uch/studspace/pkg/errors"
	"github.com/kr0uch/studspace/pkg/web"
)

type BookingService interface {
	CreateBooking(ctx context.Context, actor entities.User, in dto.CreateBooking, now time.Time) (*entities.Booking, error)
	GetBookingByID(ctx context.Context, actor entities.User, bookingID int64) (*entities.Booking, error)
	ListUserBookings(ctx context.Context, actor entities.User, in dto.ListUserBookings) ([]entities.Booking, *string, error)
	WithdrawBooking(ctx context.Context, userID, bookingID int64, now time.Time) (*entities.Booking, error)
	CancelBookingByUser(ctx context.Context, userID, bookingID int64, now time.Time) (*entities.Booking, error)
	RescheduleBooking(ctx context.Context, userID int64, in dto.RescheduleBooking, now time.Time) (*entities.Booking, error)
}

type BookingSpaceService interface {
	GetSpaceByID(ctx context.Context, spaceID int16) (*entities.Space, error)
}

type BookingUserService interface {
	GetUserByID(ctx context.Context, userID int64) (*entities.User, error)
}

type BookingController struct {
	bookingService BookingService
	spaceService   BookingSpaceService
	userService    BookingUserService
}

func NewBookingController(bookingService BookingService, spaceService BookingSpaceService, userService BookingUserService) *BookingController {
	return &BookingController{
		bookingService: bookingService,
		spaceService:   spaceService,
		userService:    userService,
	}
}

// CreateBookingHandler godoc
// @Summary Создать заявку на бронь
// @Tags Брони
// @Accept json
// @Produce json
// @Security VKAuth
// @Param request body api.CreateBookingRequest true "Слот и цель"
// @Success 201 {object} api.Booking
// @Failure 400 {object} errors.ErrorResponse
// @Failure 401 {object} errors.ErrorResponse
// @Failure 409 {object} errors.ErrorResponse
// @Failure 422 {object} errors.ErrorResponse
// @Failure 429 {object} errors.ErrorResponse
// @Failure 500 {object} errors.ErrorResponse
// @Router /bookings [post]
func (c *BookingController) CreateBookingHandler(w http.ResponseWriter, r *http.Request) error {
	user, err := web.GetUserFromCtx(r.Context())
	if err != nil {
		return err
	}

	var in api.CreateBookingRequest
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		return errors.ErrBadRequest
	}

	now := time.Now()
	booking, err := c.bookingService.CreateBooking(r.Context(), *user, dto.CreateBooking{
		SpaceID:     in.SpaceID,
		StartsAt:    in.StartsAt,
		EndsAt:      in.EndsAt,
		Purpose:     in.Purpose,
		PeopleCount: in.PeopleCount,
	}, now)
	if err != nil {
		return err
	}

	out, err := c.bookingView(r, *booking, now)
	if err != nil {
		return err
	}
	w.WriteHeader(http.StatusCreated)
	return json.NewEncoder(w).Encode(out)
}

// ListMyBookingsHandler godoc
// @Summary Мои заявки
// @Tags Брони
// @Produce json
// @Security VKAuth
// @Param tab query string false "active или history"
// @Param cursor query string false "курсор"
// @Param limit query int false "1–100, по умолчанию 20"
// @Success 200 {object} api.BookingList
// @Failure 401 {object} errors.ErrorResponse
// @Failure 422 {object} errors.ErrorResponse
// @Failure 500 {object} errors.ErrorResponse
// @Router /bookings/my [get]
func (c *BookingController) ListMyBookingsHandler(w http.ResponseWriter, r *http.Request) error {
	user, err := web.GetUserFromCtx(r.Context())
	if err != nil {
		return err
	}

	tab := entities.UserBookingTab(r.URL.Query().Get("tab"))
	limit, err := parseLimit(r.URL.Query().Get("limit"))
	if err != nil {
		return err
	}

	now := time.Now()
	bookings, nextCursor, err := c.bookingService.ListUserBookings(r.Context(), *user, dto.ListUserBookings{
		UserID: user.UserID,
		Tab:    tab,
		Now:    now,
		Cursor: r.URL.Query().Get("cursor"),
		Limit:  limit,
	})
	if err != nil {
		return err
	}

	items := make([]api.Booking, 0, len(bookings))
	for _, booking := range bookings {
		item, err := c.bookingView(r, booking, now)
		if err != nil {
			return err
		}
		items = append(items, item)
	}
	return json.NewEncoder(w).Encode(api.BookingList{Items: items, NextCursor: nextCursor})
}

// WithdrawBookingHandler godoc
// @Summary Отозвать заявку
// @Tags Брони
// @Produce json
// @Security VKAuth
// @Param id path int true "booking_id"
// @Success 200 {object} api.Booking
// @Failure 401 {object} errors.ErrorResponse
// @Failure 403 {object} errors.ErrorResponse
// @Failure 404 {object} errors.ErrorResponse
// @Failure 409 {object} errors.ErrorResponse
// @Failure 500 {object} errors.ErrorResponse
// @Router /bookings/{id}/withdraw [post]
func (c *BookingController) WithdrawBookingHandler(w http.ResponseWriter, r *http.Request) error {
	user, err := web.GetUserFromCtx(r.Context())
	if err != nil {
		return err
	}

	bookingID, err := parseInt64ID(r.PathValue("id"), "id")
	if err != nil {
		return err
	}

	now := time.Now()
	booking, err := c.bookingService.WithdrawBooking(r.Context(), user.UserID, bookingID, now)
	if err != nil {
		return err
	}

	out, err := c.bookingView(r, *booking, now)
	if err != nil {
		return err
	}
	return json.NewEncoder(w).Encode(out)
}

// CancelBookingHandler godoc
// @Summary Отменить одобренную бронь
// @Tags Брони
// @Produce json
// @Security VKAuth
// @Param id path int true "booking_id"
// @Success 200 {object} api.Booking
// @Failure 401 {object} errors.ErrorResponse
// @Failure 403 {object} errors.ErrorResponse
// @Failure 404 {object} errors.ErrorResponse
// @Failure 409 {object} errors.ErrorResponse
// @Failure 500 {object} errors.ErrorResponse
// @Router /bookings/{id}/cancel [post]
func (c *BookingController) CancelBookingHandler(w http.ResponseWriter, r *http.Request) error {
	user, err := web.GetUserFromCtx(r.Context())
	if err != nil {
		return err
	}

	bookingID, err := parseInt64ID(r.PathValue("id"), "id")
	if err != nil {
		return err
	}

	now := time.Now()
	booking, err := c.bookingService.CancelBookingByUser(r.Context(), user.UserID, bookingID, now)
	if err != nil {
		return err
	}

	out, err := c.bookingView(r, *booking, now)
	if err != nil {
		return err
	}
	return json.NewEncoder(w).Encode(out)
}

// RescheduleBookingHandler godoc
// @Summary Перенести бронь
// @Tags Брони
// @Accept json
// @Produce json
// @Security VKAuth
// @Param id path int true "booking_id"
// @Param request body api.RescheduleBookingRequest true "Новый слот"
// @Success 200 {object} api.Booking
// @Failure 400 {object} errors.ErrorResponse
// @Failure 401 {object} errors.ErrorResponse
// @Failure 403 {object} errors.ErrorResponse
// @Failure 409 {object} errors.ErrorResponse
// @Failure 422 {object} errors.ErrorResponse
// @Failure 500 {object} errors.ErrorResponse
// @Router /bookings/{id}/reschedule [post]
func (c *BookingController) RescheduleBookingHandler(w http.ResponseWriter, r *http.Request) error {
	user, err := web.GetUserFromCtx(r.Context())
	if err != nil {
		return err
	}

	bookingID, err := parseInt64ID(r.PathValue("id"), "id")
	if err != nil {
		return err
	}

	var in api.RescheduleBookingRequest
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		return errors.ErrBadRequest
	}

	now := time.Now()
	spaceID := in.SpaceID
	if spaceID < 1 {
		current, err := c.bookingService.GetBookingByID(r.Context(), *user, bookingID)
		if err != nil {
			return err
		}
		spaceID = current.SpaceID
	}

	booking, err := c.bookingService.RescheduleBooking(r.Context(), user.UserID, dto.RescheduleBooking{
		BookingID: bookingID,
		SpaceID:   spaceID,
		StartsAt:  in.StartsAt,
		EndsAt:    in.EndsAt,
	}, now)
	if err != nil {
		return err
	}

	out, err := c.bookingView(r, *booking, now)
	if err != nil {
		return err
	}
	return json.NewEncoder(w).Encode(out)
}

func (c *BookingController) bookingView(r *http.Request, booking entities.Booking, now time.Time) (api.Booking, error) {
	space, err := c.spaceService.GetSpaceByID(r.Context(), booking.SpaceID)
	if err != nil {
		space = nil
	}

	var decisionBy *entities.User
	if booking.DecidedBy.Valid {
		decisionBy, _ = c.userService.GetUserByID(r.Context(), booking.DecidedBy.Int64)
	}
	return api.NewBooking(booking, space, decisionBy, now), nil
}
