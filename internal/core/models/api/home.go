package api

import (
	"time"

	"github.com/kr0uch/studspace/internal/core/models/entities"
)

type HomeTodaySpace struct {
	SpaceID int16  `json:"space_id"`
	Name    string `json:"name"`
}

type HomeTodayItem struct {
	BookingID int64                  `json:"booking_id"`
	StartsAt  time.Time              `json:"starts_at"`
	EndsAt    time.Time              `json:"ends_at"`
	Space     HomeTodaySpace         `json:"space"`
	Title     string                 `json:"title"`
	Purpose   string                 `json:"purpose"`
	Status    entities.BookingStatus `json:"status"`
}

type HomeToday struct {
	Total int             `json:"total"`
	Items []HomeTodayItem `json:"items"`
}

type Home struct {
	Today  HomeToday `json:"today"`
	Spaces []Space   `json:"spaces"`
}
