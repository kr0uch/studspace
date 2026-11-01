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

type AdminBookingService interface {
	ListAdminBookings(ctx context.Context, actor entities.User, in dto.ListAdminBookings) ([]entities.Booking, *string, error)
	CountBookingsByStatus(ctx context.Context, actor entities.User, status entities.BookingStatus) (int, error)
	GetSlotOverlaps(ctx context.Context, actor entities.User, in dto.GetSlotOverlaps) (*dto.SlotOverlaps, error)
	ApproveBooking(ctx context.Context, actor entities.User, bookingID int64, now time.Time) (*entities.Booking, error)
	RejectBooking(ctx context.Context, actor entities.User, bookingID int64, reason entities.RejectReason, now time.Time) (*entities.Booking, error)
	CancelBookingByAdmin(ctx context.Context, actor entities.User, bookingID int64, now time.Time) (*entities.Booking, error)
	ReturnBookingToReview(ctx context.Context, actor entities.User, bookingID int64, now time.Time) (*entities.Booking, error)
}

type AdminEventService interface {
	CreateEvent(ctx context.Context, actor entities.User, in dto.CreateEvent, now time.Time) (*entities.Event, error)
	UpdateEvent(ctx context.Context, actor entities.User, eventID int64, in dto.CreateEvent, now time.Time) (*entities.Event, error)
	CancelEvent(ctx context.Context, actor entities.User, eventID int64) (*entities.Event, error)
}

type AdminSpaceService interface {
	GetSpaceByID(ctx context.Context, spaceID int16) (*entities.Space, error)
}

type AdminUserService interface {
	GetUserByID(ctx context.Context, userID int64) (*entities.User, error)
}

type AdminController struct {
	bookingService AdminBookingService
	eventService   AdminEventService
	spaceService   AdminSpaceService
	userService    AdminUserService
}

func NewAdminController(
	bookingService AdminBookingService,
	eventService AdminEventService,
	spaceService AdminSpaceService,
	userService AdminUserService,
) *AdminController {
	return &AdminController{
		bookingService: bookingService,
		eventService:   eventService,
		spaceService:   spaceService,
		userService:    userService,
	}
}

// ListAdminBookingsHandler godoc
// @Summary Админ: список заявок
// @Tags Админ
// @Produce json
// @Security VKAuth
// @Param tab query string false "new, approved, rejected"
// @Param cursor query string false "курсор"
// @Param limit query int false "1–100, по умолчанию 20"
// @Success 200 {object} api.AdminBookingList
// @Failure 401 {object} errors.ErrorResponse
// @Failure 403 {object} errors.ErrorResponse
// @Failure 422 {object} errors.ErrorResponse
// @Failure 500 {object} errors.ErrorResponse
// @Router /admin/bookings [get]
func (c *AdminController) ListAdminBookingsHandler(w http.ResponseWriter, r *http.Request) error {
	user, err := web.GetUserFromCtx(r.Context())
	if err != nil {
		return err
	}

	tab := entities.AdminBookingTab(r.URL.Query().Get("tab"))
	limit, err := parseLimit(r.URL.Query().Get("limit"))
	if err != nil {
		return err
	}

	now := time.Now()
	bookings, nextCursor, err := c.bookingService.ListAdminBookings(r.Context(), *user, dto.ListAdminBookings{
		Tab:    tab,
		Cursor: r.URL.Query().Get("cursor"),
		Limit:  limit,
	})
	if err != nil {
		return err
	}

	newCount, err := c.bookingService.CountBookingsByStatus(r.Context(), *user, entities.BookingStatusPending)
	if err != nil {
		return err
	}

	items := make([]api.AdminBooking, 0, len(bookings))
	for _, booking := range bookings {
		item, err := c.adminBookingView(r, *user, booking, now)
		if err != nil {
			return err
		}
		items = append(items, item)
	}

	return json.NewEncoder(w).Encode(api.AdminBookingList{
		Counters:   api.AdminBookingCounters{New: newCount},
		Items:      items,
		NextCursor: nextCursor,
	})
}

// ApproveBookingHandler godoc
// @Summary Админ: одобрить заявку
// @Tags Админ
// @Produce json
// @Security VKAuth
// @Param id path int true "booking_id"
// @Success 200 {object} api.AdminBooking
// @Failure 401 {object} errors.ErrorResponse
// @Failure 403 {object} errors.ErrorResponse
// @Failure 409 {object} errors.ErrorResponse
// @Failure 500 {object} errors.ErrorResponse
// @Router /admin/bookings/{id}/approve [post]
func (c *AdminController) ApproveBookingHandler(w http.ResponseWriter, r *http.Request) error {
	user, err := web.GetUserFromCtx(r.Context())
	if err != nil {
		return err
	}

	bookingID, err := parseInt64ID(r.PathValue("id"), "id")
	if err != nil {
		return err
	}

	now := time.Now()
	booking, err := c.bookingService.ApproveBooking(r.Context(), *user, bookingID, now)
	if booking == nil {
		return err
	}

	out, viewErr := c.adminBookingView(r, *user, *booking, now)
	if viewErr != nil {
		return viewErr
	}
	if encErr := json.NewEncoder(w).Encode(out); encErr != nil {
		return encErr
	}
	return err
}

// RejectBookingHandler godoc
// @Summary Админ: отклонить заявку
// @Tags Админ
// @Accept json
// @Produce json
// @Security VKAuth
// @Param id path int true "booking_id"
// @Param request body api.RejectBookingRequest true "Причина"
// @Success 200 {object} api.AdminBooking
// @Failure 400 {object} errors.ErrorResponse
// @Failure 401 {object} errors.ErrorResponse
// @Failure 403 {object} errors.ErrorResponse
// @Failure 409 {object} errors.ErrorResponse
// @Failure 500 {object} errors.ErrorResponse
// @Router /admin/bookings/{id}/reject [post]
func (c *AdminController) RejectBookingHandler(w http.ResponseWriter, r *http.Request) error {
	user, err := web.GetUserFromCtx(r.Context())
	if err != nil {
		return err
	}

	bookingID, err := parseInt64ID(r.PathValue("id"), "id")
	if err != nil {
		return err
	}

	var in api.RejectBookingRequest
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		return errors.ErrBadRequest
	}

	now := time.Now()
	booking, err := c.bookingService.RejectBooking(r.Context(), *user, bookingID, entities.RejectReason{
		Title: in.ReasonTitle,
		Text:  in.ReasonText,
	}, now)
	if booking == nil {
		return err
	}

	out, viewErr := c.adminBookingView(r, *user, *booking, now)
	if viewErr != nil {
		return viewErr
	}
	if encErr := json.NewEncoder(w).Encode(out); encErr != nil {
		return encErr
	}
	return err
}

// CancelBookingHandler godoc
// @Summary Админ: отменить бронь
// @Tags Админ
// @Produce json
// @Security VKAuth
// @Param id path int true "booking_id"
// @Success 200 {object} api.AdminBooking
// @Failure 401 {object} errors.ErrorResponse
// @Failure 403 {object} errors.ErrorResponse
// @Failure 409 {object} errors.ErrorResponse
// @Failure 500 {object} errors.ErrorResponse
// @Router /admin/bookings/{id}/cancel [post]
func (c *AdminController) CancelBookingHandler(w http.ResponseWriter, r *http.Request) error {
	user, err := web.GetUserFromCtx(r.Context())
	if err != nil {
		return err
	}

	bookingID, err := parseInt64ID(r.PathValue("id"), "id")
	if err != nil {
		return err
	}

	now := time.Now()
	booking, err := c.bookingService.CancelBookingByAdmin(r.Context(), *user, bookingID, now)
	if booking == nil {
		return err
	}

	out, viewErr := c.adminBookingView(r, *user, *booking, now)
	if viewErr != nil {
		return viewErr
	}
	if encErr := json.NewEncoder(w).Encode(out); encErr != nil {
		return encErr
	}
	return err
}

// ReturnBookingToReviewHandler godoc
// @Summary Админ: вернуть на рассмотрение
// @Tags Админ
// @Produce json
// @Security VKAuth
// @Param id path int true "booking_id"
// @Success 200 {object} api.AdminBooking
// @Failure 401 {object} errors.ErrorResponse
// @Failure 403 {object} errors.ErrorResponse
// @Failure 409 {object} errors.ErrorResponse
// @Failure 500 {object} errors.ErrorResponse
// @Router /admin/bookings/{id}/return-to-review [post]
func (c *AdminController) ReturnBookingToReviewHandler(w http.ResponseWriter, r *http.Request) error {
	user, err := web.GetUserFromCtx(r.Context())
	if err != nil {
		return err
	}

	bookingID, err := parseInt64ID(r.PathValue("id"), "id")
	if err != nil {
		return err
	}

	now := time.Now()
	booking, err := c.bookingService.ReturnBookingToReview(r.Context(), *user, bookingID, now)
	if err != nil {
		return err
	}

	out, err := c.adminBookingView(r, *user, *booking, now)
	if err != nil {
		return err
	}
	return json.NewEncoder(w).Encode(out)
}

// CreateEventHandler godoc
// @Summary Админ: создать мероприятие
// @Tags Админ
// @Accept json
// @Produce json
// @Security VKAuth
// @Param request body api.CreateEventRequest true "Мероприятие"
// @Success 201 {object} api.Event
// @Failure 400 {object} errors.ErrorResponse
// @Failure 401 {object} errors.ErrorResponse
// @Failure 403 {object} errors.ErrorResponse
// @Failure 409 {object} errors.ErrorResponse
// @Failure 422 {object} errors.ErrorResponse
// @Failure 500 {object} errors.ErrorResponse
// @Router /admin/events [post]
func (c *AdminController) CreateEventHandler(w http.ResponseWriter, r *http.Request) error {
	user, err := web.GetUserFromCtx(r.Context())
	if err != nil {
		return err
	}

	var in api.CreateEventRequest
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		return errors.ErrBadRequest
	}

	event, err := c.eventService.CreateEvent(r.Context(), *user, dto.CreateEvent{
		Title:       in.Title,
		Organizer:   in.Organizer,
		Category:    in.Category,
		SpaceID:     in.SpaceID,
		StartsAt:    in.StartsAt,
		EndsAt:      in.EndsAt,
		Capacity:    in.Capacity,
		Description: in.Description,
	}, time.Now())
	if err != nil {
		return err
	}

	spaceName := ""
	if space, err := c.spaceService.GetSpaceByID(r.Context(), event.SpaceID); err == nil {
		spaceName = space.Name
	}

	w.WriteHeader(http.StatusCreated)
	return json.NewEncoder(w).Encode(api.NewEvent(*event, spaceName))
}

// UpdateEventHandler godoc
// @Summary Админ: изменить мероприятие
// @Tags Админ
// @Accept json
// @Produce json
// @Security VKAuth
// @Param id path int true "event_id"
// @Param request body api.CreateEventRequest true "То же тело, что при создании"
// @Success 200 {object} api.Event
// @Failure 400 {object} errors.ErrorResponse
// @Failure 401 {object} errors.ErrorResponse
// @Failure 403 {object} errors.ErrorResponse
// @Failure 409 {object} errors.ErrorResponse
// @Failure 422 {object} errors.ErrorResponse
// @Failure 500 {object} errors.ErrorResponse
// @Router /admin/events/{id} [patch]
func (c *AdminController) UpdateEventHandler(w http.ResponseWriter, r *http.Request) error {
	user, err := web.GetUserFromCtx(r.Context())
	if err != nil {
		return err
	}

	eventID, err := parseInt64ID(r.PathValue("id"), "id")
	if err != nil {
		return err
	}

	var in api.CreateEventRequest
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		return errors.ErrBadRequest
	}

	event, err := c.eventService.UpdateEvent(r.Context(), *user, eventID, dto.CreateEvent{
		Title:       in.Title,
		Organizer:   in.Organizer,
		Category:    in.Category,
		SpaceID:     in.SpaceID,
		StartsAt:    in.StartsAt,
		EndsAt:      in.EndsAt,
		Capacity:    in.Capacity,
		Description: in.Description,
	}, time.Now())
	if err != nil {
		return err
	}

	spaceName := ""
	if space, err := c.spaceService.GetSpaceByID(r.Context(), event.SpaceID); err == nil {
		spaceName = space.Name
	}
	return json.NewEncoder(w).Encode(api.NewEvent(*event, spaceName))
}

// CancelEventHandler godoc
// @Summary Админ: отменить мероприятие
// @Tags Админ
// @Produce json
// @Security VKAuth
// @Param id path int true "event_id"
// @Success 200 {object} api.Event
// @Failure 401 {object} errors.ErrorResponse
// @Failure 403 {object} errors.ErrorResponse
// @Failure 409 {object} errors.ErrorResponse
// @Failure 500 {object} errors.ErrorResponse
// @Router /admin/events/{id}/cancel [post]
func (c *AdminController) CancelEventHandler(w http.ResponseWriter, r *http.Request) error {
	user, err := web.GetUserFromCtx(r.Context())
	if err != nil {
		return err
	}

	eventID, err := parseInt64ID(r.PathValue("id"), "id")
	if err != nil {
		return err
	}

	event, err := c.eventService.CancelEvent(r.Context(), *user, eventID)
	if err != nil {
		return err
	}

	spaceName := ""
	if space, err := c.spaceService.GetSpaceByID(r.Context(), event.SpaceID); err == nil {
		spaceName = space.Name
	}
	return json.NewEncoder(w).Encode(api.NewEvent(*event, spaceName))
}

func (c *AdminController) adminBookingView(
	r *http.Request,
	actor entities.User,
	booking entities.Booking,
	now time.Time,
) (api.AdminBooking, error) {
	space, err := c.spaceService.GetSpaceByID(r.Context(), booking.SpaceID)
	if err != nil {
		space = nil
	}

	owner, _ := c.userService.GetUserByID(r.Context(), booking.UserID)
	var decisionBy *entities.User
	if booking.DecidedBy.Valid {
		decisionBy, _ = c.userService.GetUserByID(r.Context(), booking.DecidedBy.Int64)
	}

	overlaps, err := c.bookingService.GetSlotOverlaps(r.Context(), actor, dto.GetSlotOverlaps{
		SpaceID:          booking.SpaceID,
		StartsAt:         booking.StartsAt,
		EndsAt:           booking.EndsAt,
		ExcludeBookingID: &booking.BookingID,
	})
	if err != nil {
		return api.AdminBooking{}, err
	}

	conflicts := make([]api.AdminBookingConflict, 0, len(overlaps.Pending)+len(overlaps.Approved))
	for _, other := range append(append([]entities.Booking{}, overlaps.Pending...), overlaps.Approved...) {
		conflicts = append(conflicts, api.AdminBookingConflict{
			BookingID: other.BookingID,
			OwnerName: other.FirstName + " " + other.LastName,
			StartsAt:  other.StartsAt,
			EndsAt:    other.EndsAt,
			Status:    other.CurrentStatus(now),
		})
	}

	return api.AdminBooking{
		Booking:   api.NewBooking(booking, space, decisionBy, now),
		User:      owner,
		Conflicts: conflicts,
	}, nil
}
