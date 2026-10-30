package utils

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
	"time"
)

func EncodeCursor(at time.Time, id int64) string {
	raw := at.UTC().Format(time.RFC3339Nano) + "|" + strconv.FormatInt(id, 10)
	return base64.RawURLEncoding.EncodeToString([]byte(raw))
}

func DecodeCursor(raw string) (time.Time, int64, error) {
	b, err := base64.RawURLEncoding.DecodeString(raw)
	if err != nil {
		return time.Time{}, 0, err
	}
	atRaw, idRaw, ok := strings.Cut(string(b), "|")
	if !ok {
		return time.Time{}, 0, fmt.Errorf("cursor")
	}
	at, err := time.Parse(time.RFC3339Nano, atRaw)
	if err != nil {
		return time.Time{}, 0, err
	}
	id, err := strconv.ParseInt(idRaw, 10, 64)
	if err != nil || id < 1 {
		return time.Time{}, 0, fmt.Errorf("cursor")
	}
	return at, id, nil
}
