package routers

import (
	"net/http"

	"github.com/kr0uch/studspace/internal/transport/api/handlers"
	"github.com/kr0uch/studspace/pkg/web"
)

func NewAdminRouter(
	bookingService handlers.AdminBookingService,
	eventService handlers.AdminEventService,
	spaceService handlers.AdminSpaceService,
	userService handlers.AdminUserService,
) *http.ServeMux {
	controller := handlers.NewAdminController(bookingService, eventService, spaceService, userService)
	router := http.NewServeMux()

	router.HandleFunc("GET /admin/bookings", web.HandleWithError(controller.ListAdminBookingsHandler))
	router.HandleFunc("POST /admin/bookings/{id}/approve", web.HandleWithError(controller.ApproveBookingHandler))
	router.HandleFunc("POST /admin/bookings/{id}/reject", web.HandleWithError(controller.RejectBookingHandler))
	router.HandleFunc("POST /admin/bookings/{id}/cancel", web.HandleWithError(controller.CancelBookingHandler))
	router.HandleFunc("POST /admin/bookings/{id}/return-to-review", web.HandleWithError(controller.ReturnBookingToReviewHandler))
	router.HandleFunc("POST /admin/events", web.HandleWithError(controller.CreateEventHandler))
	router.HandleFunc("PATCH /admin/events/{id}", web.HandleWithError(controller.UpdateEventHandler))
	router.HandleFunc("POST /admin/events/{id}/cancel", web.HandleWithError(controller.CancelEventHandler))

	return router
}
