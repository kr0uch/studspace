package routers

import (
	"net/http"

	"github.com/kr0uch/studspace/internal/transport/api/handlers"
	"github.com/kr0uch/studspace/pkg/web"
)

func NewHomeRouter(bookingService handlers.HomeBookingService, eventService handlers.HomeEventService, spaceService handlers.HomeSpaceService) *http.ServeMux {
	controller := handlers.NewHomeController(bookingService, eventService, spaceService)
	router := http.NewServeMux()

	router.HandleFunc("GET /home", web.HandleWithError(controller.GetHomeHandler))

	return router
}
