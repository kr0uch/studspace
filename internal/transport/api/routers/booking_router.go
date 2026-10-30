package routers

import (
	"net/http"

	"github.com/kr0uch/studspace/internal/transport/api/handlers"
	"github.com/kr0uch/studspace/pkg/web"
)

func NewBookingRouter(bookingService handlers.BookingService, spaceService handlers.BookingSpaceService, userService handlers.BookingUserService) *http.ServeMux {
	controller := handlers.NewBookingController(bookingService, spaceService, userService)
	router := http.NewServeMux()

	router.HandleFunc("POST /bookings", web.HandleWithError(controller.CreateBookingHandler))
	router.HandleFunc("GET /bookings/my", web.HandleWithError(controller.ListMyBookingsHandler))
	router.HandleFunc("POST /bookings/{id}/withdraw", web.HandleWithError(controller.WithdrawBookingHandler))
	router.HandleFunc("POST /bookings/{id}/cancel", web.HandleWithError(controller.CancelBookingHandler))
	router.HandleFunc("POST /bookings/{id}/reschedule", web.HandleWithError(controller.RescheduleBookingHandler))

	return router
}
