package service

import (
	"fmt"
	"time"

	"github.com/kr0uch/studspace/internal/core/models/entities"
	"github.com/kr0uch/studspace/pkg/utils"
)

const (
	msgBookingApproved       = "Заявка на %s %s одобрена."
	msgBookingRejected       = "Заявка на %s %s отклонена."
	msgBookingRejectedReason = "Заявка на %s %s отклонена: %s"
	msgBookingCancelled      = "Бронь %s %s отменена администратором."
	msgWaitlistFreed         = "Освободилось место на «%s». Запишитесь в приложении."
)

func bookingApprovedMessage(spaceName string, booking entities.Booking) string {
	return fmt.Sprintf(msgBookingApproved, spaceName, formatSlot(booking.StartsAt, booking.EndsAt))
}

func bookingRejectedMessage(spaceName string, booking entities.Booking, reason string) string {
	slot := formatSlot(booking.StartsAt, booking.EndsAt)
	if reason == "" {
		return fmt.Sprintf(msgBookingRejected, spaceName, slot)
	}
	return fmt.Sprintf(msgBookingRejectedReason, spaceName, slot, reason)
}

func bookingCancelledMessage(spaceName string, booking entities.Booking) string {
	return fmt.Sprintf(msgBookingCancelled, spaceName, formatSlot(booking.StartsAt, booking.EndsAt))
}

func waitlistFreedMessage(title string) string {
	return fmt.Sprintf(msgWaitlistFreed, title)
}

func allowedVKUserIDs(users []entities.User) []int64 {
	ids := make([]int64, 0, len(users))
	for _, user := range users {
		if user.NotificationsAllowed && user.VKUserID > 0 {
			ids = append(ids, user.VKUserID)
		}
	}
	return ids
}

func formatSlot(startsAt, endsAt time.Time) string {
	start := startsAt.In(utils.Moscow)
	end := endsAt.In(utils.Moscow)
	return fmt.Sprintf("%02d.%02d %02d:%02d–%02d:%02d",
		start.Day(), int(start.Month()), start.Hour(), start.Minute(),
		end.Hour(), end.Minute(),
	)
}
