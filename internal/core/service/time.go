package service

import (
	"time"

	"github.com/kr0uch/studspace/internal/core/models/entities"
	"github.com/kr0uch/studspace/pkg/errors"
)

var moscow *time.Location

func init() {
	loc, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		loc = time.FixedZone("MSK", 3*3600)
	}
	moscow = loc
}

func dayRange(day time.Time) (time.Time, time.Time) {
	start := time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, moscow)
	return start, start.Add(24 * time.Hour)
}

func monthRange(month time.Time) (time.Time, time.Time) {
	start := time.Date(month.Year(), month.Month(), 1, 0, 0, 0, 0, moscow)
	return start, start.AddDate(0, 1, 0)
}

func clockOf(t time.Time) time.Duration {
	local := t.In(moscow)
	return time.Duration(local.Hour())*time.Hour +
		time.Duration(local.Minute())*time.Minute +
		time.Duration(local.Second())*time.Second
}

func atClock(day time.Time, clock time.Duration) time.Time {
	start := time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, moscow)
	return start.Add(clock)
}

func validateSpaceTimeWindow(space entities.Space, startsAt, endsAt, now time.Time) error {
	fields := map[string]string{}
	if !startsAt.After(now) {
		fields["starts_at"] = "время должно быть в будущем"
	}
	if !endsAt.After(startsAt) {
		fields["ends_at"] = "окончание должно быть позже начала"
	}

	startDay := startsAt.In(moscow)
	endDay := endsAt.In(moscow)
	if startDay.Year() != endDay.Year() || startDay.YearDay() != endDay.YearDay() {
		fields["ends_at"] = "интервал должен быть в пределах одного дня"
	}

	startClock := clockOf(startsAt)
	endClock := clockOf(endsAt)
	step := time.Duration(space.SlotStepMin) * time.Minute
	if startClock < space.OpenTime || startClock >= space.CloseTime {
		fields["starts_at"] = "вне рабочих часов помещения"
	} else if step > 0 && (startClock-space.OpenTime)%step != 0 {
		fields["starts_at"] = "время должно быть кратно шагу слота"
	}
	if endClock <= space.OpenTime || endClock > space.CloseTime {
		fields["ends_at"] = "вне рабочих часов помещения"
	} else if step > 0 && endClock != space.CloseTime && (endClock-space.OpenTime)%step != 0 {
		fields["ends_at"] = "время должно быть кратно шагу слота"
	}

	if len(fields) > 0 {
		return errors.Validation(fields)
	}
	return nil
}
