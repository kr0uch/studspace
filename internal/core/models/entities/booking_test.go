package entities

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBooking_CurrentStatus(t *testing.T) {
	now := time.Date(2026, 11, 7, 12, 0, 0, 0, time.UTC)
	booking := Booking{
		Status:   BookingStatusPending,
		StartsAt: now.Add(time.Hour),
		EndsAt:   now.Add(2 * time.Hour),
	}

	t.Run("pending in the future stays pending", func(t *testing.T) {
		assert.Equal(t, BookingStatusPending, booking.CurrentStatus(now))
		assert.True(t, booking.IsActive(now))
		assert.Equal(t, []BookingAction{BookingActionWithdraw}, booking.StudentActions(now))
	})

	t.Run("pending that already started is expired", func(t *testing.T) {
		started := booking
		started.StartsAt = now
		assert.Equal(t, BookingStatusExpired, started.CurrentStatus(now))
		assert.False(t, started.IsActive(now))
		assert.Equal(t, []BookingAction{BookingActionRebook}, started.StudentActions(now))
	})

	t.Run("approved still running stays approved", func(t *testing.T) {
		running := booking
		running.Status = BookingStatusApproved
		assert.Equal(t, BookingStatusApproved, running.CurrentStatus(now))
		assert.Equal(t, []BookingAction{BookingActionReschedule, BookingActionCancel}, running.StudentActions(now))
	})

	t.Run("approved after the end is completed", func(t *testing.T) {
		done := booking
		done.Status = BookingStatusApproved
		done.EndsAt = now
		assert.Equal(t, BookingStatusCompleted, done.CurrentStatus(now))
		assert.Equal(t, []BookingAction{BookingActionRebook}, done.StudentActions(now))
	})
}

func TestBookingStatus_CanTransitionTo(t *testing.T) {
	assert.True(t, BookingStatusPending.CanTransitionTo(BookingStatusApproved))
	assert.False(t, BookingStatusExpired.CanTransitionTo(BookingStatusApproved))
	assert.False(t, BookingStatusCompleted.CanTransitionTo(BookingStatusPending))
	assert.True(t, BookingStatusRejected.CanTransitionTo(BookingStatusPending))
}
