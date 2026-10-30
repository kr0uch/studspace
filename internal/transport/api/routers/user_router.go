package routers

import (
	"net/http"

	"github.com/kr0uch/studspace/internal/transport/api/handlers"
	"github.com/kr0uch/studspace/pkg/web"
)

func NewUserRouter(userService handlers.UserService) *http.ServeMux {
	controller := handlers.NewUserController(userService)
	router := http.NewServeMux()

	router.HandleFunc("GET /me", web.HandleWithError(controller.GetMeHandler))
	router.HandleFunc("PATCH /me/notifications", web.HandleWithError(controller.PatchNotificationsHandler))

	return router
}
