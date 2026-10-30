package utils

import (
	"fmt"
	"time"
)

var Moscow *time.Location

func init() {
	loc, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		loc = time.FixedZone("MSK", 3*3600)
	}
	Moscow = loc
}

func DayRange(day time.Time) (time.Time, time.Time) {
	start := time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, Moscow)
	return start, start.Add(24 * time.Hour)
}

func MonthRange(month time.Time) (time.Time, time.Time) {
	start := time.Date(month.Year(), month.Month(), 1, 0, 0, 0, 0, Moscow)
	return start, start.AddDate(0, 1, 0)
}

func ClockOf(t time.Time) time.Duration {
	local := t.In(Moscow)
	return time.Duration(local.Hour())*time.Hour +
		time.Duration(local.Minute())*time.Minute +
		time.Duration(local.Second())*time.Second
}

func AtClock(day time.Time, clock time.Duration) time.Time {
	start := time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, Moscow)
	return start.Add(clock)
}

func FormatClock(d time.Duration) string {
	if d < 0 {
		d = 0
	}
	minutes := int(d.Minutes())
	return fmt.Sprintf("%02d:%02d", minutes/60, minutes%60)
}
