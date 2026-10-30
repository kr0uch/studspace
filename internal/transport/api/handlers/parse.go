package handlers

import (
	"strconv"
	"time"

	"github.com/kr0uch/studspace/internal/core/models/dto"
	"github.com/kr0uch/studspace/pkg/errors"
	"github.com/kr0uch/studspace/pkg/utils"
)

func parseLimit(raw string) (int, error) {
	if raw == "" {
		return dto.DefaultListLimit, nil
	}
	n, err := strconv.Atoi(raw)
	if err != nil || n < 1 || n > dto.MaxListLimit {
		return 0, errors.Validation(map[string]string{"limit": "от 1 до 100"})
	}
	return n, nil
}

func parseInt64ID(raw, field string) (int64, error) {
	id, err := strconv.ParseInt(raw, 10, 64)
	if err != nil || id < 1 {
		return 0, errors.Validation(map[string]string{field: "некорректный идентификатор"})
	}
	return id, nil
}

func parseDateField(raw, field string) (time.Time, error) {
	if raw == "" {
		return time.Time{}, errors.Validation(map[string]string{field: "обязательное поле"})
	}
	date, err := time.ParseInLocation("2006-01-02", raw, utils.Moscow)
	if err != nil {
		return time.Time{}, errors.Validation(map[string]string{field: "ожидается YYYY-MM-DD"})
	}
	return date, nil
}

func parseInclusiveDays(fromRaw, toRaw string) (time.Time, time.Time, error) {
	from, err := parseDateField(fromRaw, "from")
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	to, err := parseDateField(toRaw, "to")
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	if to.Before(from) {
		return time.Time{}, time.Time{}, errors.Validation(map[string]string{"to": "не раньше from"})
	}
	fromStart, _ := utils.DayRange(from)
	_, toEnd := utils.DayRange(to)
	return fromStart, toEnd, nil
}

func parseMonth(raw string) (time.Time, error) {
	if raw == "" {
		return time.Time{}, errors.Validation(map[string]string{"month": "обязательное поле"})
	}
	month, err := time.ParseInLocation("2006-01", raw, utils.Moscow)
	if err != nil {
		return time.Time{}, errors.Validation(map[string]string{"month": "ожидается YYYY-MM"})
	}
	return month, nil
}
