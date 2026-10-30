package routers

import (
	"net/http"

	"github.com/kr0uch/studspace/internal/transport/api/handlers"
	"github.com/kr0uch/studspace/pkg/web"
)

func NewScheduleRouter(spaceService handlers.ScheduleSpaceService, bookingService handlers.ScheduleBookingService, eventService handlers.ScheduleEventService) *http.ServeMux {
	controller := handlers.NewScheduleController(spaceService, bookingService, eventService)
	router := http.NewServeMux()

	router.HandleFunc("GET /schedule", web.HandleWithError(controller.GetScheduleHandler))
	router.HandleFunc("GET /calendar/busy-days", web.HandleWithError(controller.GetBusyDaysHandler))

	return router
}
