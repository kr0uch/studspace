package routers

import (
	"net/http"

	"github.com/kr0uch/studspace/internal/transport/api/handlers"
	"github.com/kr0uch/studspace/pkg/web"
)

func NewEventRouter(eventService handlers.EventService, spaceService handlers.EventSpaceService) *http.ServeMux {
	controller := handlers.NewEventController(eventService, spaceService)
	router := http.NewServeMux()

	router.HandleFunc("GET /events", web.HandleWithError(controller.ListEventsHandler))
	router.HandleFunc("POST /events/{id}/registration", web.HandleWithError(controller.RegisterForEventHandler))
	router.HandleFunc("DELETE /events/{id}/registration", web.HandleWithError(controller.UnregisterFromEventHandler))
	router.HandleFunc("POST /events/{id}/waitlist", web.HandleWithError(controller.AddWaitlistHandler))
	router.HandleFunc("DELETE /events/{id}/waitlist", web.HandleWithError(controller.RemoveWaitlistHandler))

	return router
}
