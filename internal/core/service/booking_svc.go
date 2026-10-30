package service

import (
	"context"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/kr0uch/studspace/internal/core/models/dto"
	"github.com/kr0uch/studspace/internal/core/models/entities"
	"github.com/kr0uch/studspace/internal/core/repository/postgres"
	"github.com/kr0uch/studspace/pkg/errors"
	"github.com/kr0uch/studspace/pkg/utils"
)

const maxPendingBookings = 3

type BookingService struct {
	bookings *postgres.BookingRepository
	events   *postgres.EventRepository
	spaces   *postgres.SpaceRepository
}

func NewBookingService(
	bookings *postgres.BookingRepository,
	events *postgres.EventRepository,
	spaces *postgres.SpaceRepository,
) *BookingService {
	return &BookingService{bookings: bookings, events: events, spaces: spaces}
}

func (s *BookingService) CreateBooking(ctx context.Context, actor entities.User, in dto.CreateBooking, now time.Time) (*entities.Booking, error) {
	space, err := activeSpace(ctx, s.spaces, in.SpaceID)
	if err != nil {
		return nil, err
	}
	in.Purpose = strings.TrimSpace(in.Purpose)
	if err := s.validateBookingInput(*space, in.StartsAt, in.EndsAt, in.Purpose, in.PeopleCount, now); err != nil {
		return nil, err
	}

	pending, err := s.bookings.CountPendingBookingsByUser(ctx, actor.UserID)
	if err != nil {
		return nil, err
	}
	if pending >= maxPendingBookings {
		return nil, errors.ErrTooManyPending
	}

	overlaps, err := loadSlotOverlaps(ctx, s.bookings, s.events, in.SpaceID, in.StartsAt, in.EndsAt, nil, nil)
	if err != nil {
		return nil, err
	}
	if overlaps.HardConflict() {
		return nil, errors.ErrSlotConflict
	}

	return s.bookings.CreateBooking(ctx, entities.Booking{
		UserID:      actor.UserID,
		SpaceID:     in.SpaceID,
		FirstName:   actor.FirstName,
		LastName:    actor.LastName,
		StartsAt:    in.StartsAt,
		EndsAt:      in.EndsAt,
		Purpose:     in.Purpose,
		PeopleCount: in.PeopleCount,
	})
}

func (s *BookingService) GetBookingByID(ctx context.Context, actor entities.User, bookingID int64) (*entities.Booking, error) {
	booking, err := s.bookings.GetBookingByID(ctx, bookingID)
	if err != nil {
		return nil, err
	}
	if booking.UserID != actor.UserID && !actor.IsAdmin() {
		return nil, errors.ErrForbidden
	}
	return booking, nil
}

func (s *BookingService) ListUserBookings(ctx context.Context, actor entities.User, in dto.ListUserBookings) ([]entities.Booking, *string, error) {
	if actor.UserID != in.UserID && !actor.IsAdmin() {
		return nil, nil, errors.ErrForbidden
	}
	switch in.Tab {
	case entities.UserBookingTabActive, entities.UserBookingTabHistory:
	default:
		return nil, nil, errors.Validation(map[string]string{"tab": "неизвестная вкладка"})
	}
	cursorAt, cursorID, limit, err := listCursor(in.Cursor, in.Limit)
	if err != nil {
		return nil, nil, err
	}
	rows, err := s.bookings.ListUserBookings(ctx, in.UserID, in.Tab, in.Now, cursorAt, cursorID, fetchLimit(limit))
	if err != nil {
		return nil, nil, err
	}
	return clipBookingPage(rows, limit, func(booking entities.Booking) (time.Time, int64) {
		return booking.StartsAt, booking.BookingID
	})
}

func (s *BookingService) ListAdminBookings(ctx context.Context, actor entities.User, in dto.ListAdminBookings) ([]entities.Booking, *string, error) {
	if err := requireAdmin(actor); err != nil {
		return nil, nil, err
	}
	switch in.Tab {
	case entities.AdminBookingTabNew, entities.AdminBookingTabApproved, entities.AdminBookingTabRejected:
	default:
		return nil, nil, errors.Validation(map[string]string{"tab": "неизвестная вкладка"})
	}
	cursorAt, cursorID, limit, err := listCursor(in.Cursor, in.Limit)
	if err != nil {
		return nil, nil, err
	}
	rows, err := s.bookings.ListAdminBookings(ctx, in.Tab, cursorAt, cursorID, fetchLimit(limit))
	if err != nil {
		return nil, nil, err
	}
	key := func(booking entities.Booking) (time.Time, int64) {
		return booking.CreatedAt, booking.BookingID
	}
	if in.Tab == entities.AdminBookingTabApproved {
		key = func(booking entities.Booking) (time.Time, int64) {
			return booking.StartsAt, booking.BookingID
		}
	}
	return clipBookingPage(rows, limit, key)
}

func (s *BookingService) CountBookingsByStatus(ctx context.Context, actor entities.User, status entities.BookingStatus) (int, error) {
	if err := requireAdmin(actor); err != nil {
		return 0, err
	}
	return s.bookings.CountBookingsByStatus(ctx, status)
}

func (s *BookingService) GetSlotOverlaps(ctx context.Context, actor entities.User, in dto.GetSlotOverlaps) (*dto.SlotOverlaps, error) {
	if err := requireAdmin(actor); err != nil {
		return nil, err
	}
	return loadSlotOverlaps(ctx, s.bookings, s.events, in.SpaceID, in.StartsAt, in.EndsAt, in.ExcludeBookingID, nil)
}

func (s *BookingService) ListBusyBookingDays(ctx context.Context, month time.Time) ([]time.Time, error) {
	if month.IsZero() {
		return nil, errors.Validation(map[string]string{"month": "обязательное поле"})
	}
	from, to := utils.MonthRange(month)
	return s.bookings.ListBusyBookingDays(ctx, from, to)
}

func (s *BookingService) ListBookingsInRange(ctx context.Context, from, to time.Time) ([]entities.Booking, error) {
	if from.IsZero() || to.IsZero() || !to.After(from) {
		return nil, errors.Validation(map[string]string{"from": "укажите корректный интервал"})
	}
	return s.bookings.ListBookingsInRange(ctx, from, to, nil)
}

func (s *BookingService) WithdrawBooking(ctx context.Context, userID, bookingID int64, now time.Time) (*entities.Booking, error) {
	if _, err := s.ownBooking(ctx, userID, bookingID); err != nil {
		return nil, err
	}
	return s.bookings.WithdrawBooking(ctx, bookingID, userID, now)
}

func (s *BookingService) CancelBookingByUser(ctx context.Context, userID, bookingID int64, now time.Time) (*entities.Booking, error) {
	if _, err := s.ownBooking(ctx, userID, bookingID); err != nil {
		return nil, err
	}
	return s.bookings.CancelBooking(ctx, bookingID, userID, now, entities.CancelActorUser)
}

func (s *BookingService) RescheduleBooking(ctx context.Context, userID int64, in dto.RescheduleBooking, now time.Time) (*entities.Booking, error) {
	current, err := s.ownBooking(ctx, userID, in.BookingID)
	if err != nil {
		return nil, err
	}
	space, err := activeSpace(ctx, s.spaces, in.SpaceID)
	if err != nil {
		return nil, err
	}
	if err := s.validateBookingInput(*space, in.StartsAt, in.EndsAt, current.Purpose, current.PeopleCount, now); err != nil {
		return nil, err
	}

	exclude := current.BookingID
	overlaps, err := loadSlotOverlaps(ctx, s.bookings, s.events, in.SpaceID, in.StartsAt, in.EndsAt, &exclude, nil)
	if err != nil {
		return nil, err
	}
	if overlaps.HardConflict() {
		return nil, errors.ErrSlotConflict
	}

	return s.bookings.RescheduleBooking(ctx, in.BookingID, userID, now, in.SpaceID, in.StartsAt, in.EndsAt)
}

func (s *BookingService) ApproveBooking(ctx context.Context, actor entities.User, bookingID int64, now time.Time) (*entities.Booking, error) {
	if err := requireAdmin(actor); err != nil {
		return nil, err
	}
	booking, err := s.bookings.GetBookingByID(ctx, bookingID)
	if err != nil {
		return nil, err
	}
	exclude := booking.BookingID
	overlaps, err := loadSlotOverlaps(ctx, s.bookings, s.events, booking.SpaceID, booking.StartsAt, booking.EndsAt, &exclude, nil)
	if err != nil {
		return nil, err
	}
	if overlaps.HardConflict() {
		return nil, errors.ErrSlotConflict
	}
	return s.bookings.ApproveBooking(ctx, bookingID, actor.UserID, now)
}

func (s *BookingService) RejectBooking(ctx context.Context, actor entities.User, bookingID int64, reason entities.RejectReason, now time.Time) (*entities.Booking, error) {
	if err := requireAdmin(actor); err != nil {
		return nil, err
	}
	if strings.TrimSpace(reason.Text) == "" {
		return nil, errors.Validation(map[string]string{"reason_text": "укажите причину"})
	}
	return s.bookings.RejectBooking(ctx, bookingID, actor.UserID, now, reason)
}

func (s *BookingService) CancelBookingByAdmin(ctx context.Context, actor entities.User, bookingID int64, now time.Time) (*entities.Booking, error) {
	if err := requireAdmin(actor); err != nil {
		return nil, err
	}
	return s.bookings.CancelBooking(ctx, bookingID, actor.UserID, now, entities.CancelActorAdmin)
}

func (s *BookingService) ReturnBookingToReview(ctx context.Context, actor entities.User, bookingID int64, now time.Time) (*entities.Booking, error) {
	if err := requireAdmin(actor); err != nil {
		return nil, err
	}
	return s.bookings.ReturnBookingToReview(ctx, bookingID, actor.UserID, now)
}

func (s *BookingService) ownBooking(ctx context.Context, userID, bookingID int64) (*entities.Booking, error) {
	booking, err := s.bookings.GetBookingByID(ctx, bookingID)
	if err != nil {
		return nil, err
	}
	if booking.UserID != userID {
		return nil, errors.ErrForbidden
	}
	return booking, nil
}

func (s *BookingService) validateBookingInput(space entities.Space, startsAt, endsAt time.Time, purpose string, peopleCount int, now time.Time) error {
	if err := validateSpaceTimeWindow(space, startsAt, endsAt, now); err != nil {
		return err
	}
	fields := map[string]string{}
	if n := utf8.RuneCountInString(purpose); n < 1 || n > 150 {
		fields["purpose"] = "от 1 до 150 символов"
	}
	if peopleCount < 1 || peopleCount > space.Capacity {
		fields["people_count"] = "не больше вместимости помещения"
	}
	if len(fields) > 0 {
		return errors.Validation(fields)
	}
	return nil
}

func listCursor(raw string, limit int) (*time.Time, int64, int, error) {
	if limit < 0 || limit > dto.MaxListLimit {
		return nil, 0, 0, errors.Validation(map[string]string{"limit": "от 1 до 100"})
	}
	if raw == "" {
		return nil, 0, limit, nil
	}
	at, id, err := utils.DecodeCursor(raw)
	if err != nil {
		return nil, 0, 0, errors.Validation(map[string]string{"cursor": "некорректный курсор"})
	}
	return &at, id, limit, nil
}

func fetchLimit(limit int) int {
	if limit < 1 {
		return 0
	}
	return limit + 1
}

func clipBookingPage(items []entities.Booking, limit int, key func(entities.Booking) (time.Time, int64)) ([]entities.Booking, *string, error) {
	if limit < 1 || len(items) <= limit {
		return items, nil, nil
	}
	items = items[:limit]
	at, id := key(items[len(items)-1])
	cursor := utils.EncodeCursor(at, id)
	return items, &cursor, nil
}

func loadSlotOverlaps(ctx context.Context, bookings *postgres.BookingRepository, events *postgres.EventRepository, spaceID int16, startsAt, endsAt time.Time, excludeBookingID, excludeEventID *int64) (*dto.SlotOverlaps, error) {
	overlappingBookings, err := bookings.ListOverlappingBookings(ctx, spaceID, startsAt, endsAt, excludeBookingID)
	if err != nil {
		return nil, err
	}
	overlappingEvents, err := events.ListOverlappingEvents(ctx, spaceID, startsAt, endsAt, excludeEventID)
	if err != nil {
		return nil, err
	}

	result := &dto.SlotOverlaps{Events: overlappingEvents}
	for _, booking := range overlappingBookings {
		if booking.Status == entities.BookingStatusApproved {
			result.Approved = append(result.Approved, booking)
			continue
		}
		result.Pending = append(result.Pending, booking)
	}
	return result, nil
}
