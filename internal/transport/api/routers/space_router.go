package routers

import (
	"net/http"

	"github.com/kr0uch/studspace/internal/transport/api/handlers"
	"github.com/kr0uch/studspace/pkg/web"
)

func NewSpaceRouter(spaceService handlers.SpaceService) *http.ServeMux {
	controller := handlers.NewSpaceController(spaceService)
	router := http.NewServeMux()

	router.HandleFunc("GET /spaces", web.HandleWithError(controller.ListSpacesHandler))
	router.HandleFunc("GET /spaces/{id}/availability", web.HandleWithError(controller.GetAvailabilityHandler))

	return router
}
