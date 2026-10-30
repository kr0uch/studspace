package service

import (
	"time"

	"github.com/kr0uch/studspace/internal/core/models/entities"
	"github.com/kr0uch/studspace/pkg/errors"
	"github.com/kr0uch/studspace/pkg/utils"
)

func validateSpaceTimeWindow(space entities.Space, startsAt, endsAt, now time.Time) error {
	fields := map[string]string{}
	if !startsAt.After(now) {
		fields["starts_at"] = "время должно быть в будущем"
	}
	if !endsAt.After(startsAt) {
		fields["ends_at"] = "окончание должно быть позже начала"
	}

	startDay := startsAt.In(utils.Moscow)
	endDay := endsAt.In(utils.Moscow)
	if startDay.Year() != endDay.Year() || startDay.YearDay() != endDay.YearDay() {
		fields["ends_at"] = "интервал должен быть в пределах одного дня"
	}

	startClock := utils.ClockOf(startsAt)
	endClock := utils.ClockOf(endsAt)
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
