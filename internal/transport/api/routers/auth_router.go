package routers

import (
	"net/http"

	"github.com/kr0uch/studspace/internal/transport/api/handlers"
	"github.com/kr0uch/studspace/pkg/web"
)

func NewAuthRouter(userService handlers.UserService) *http.ServeMux {
	controller := handlers.NewUserController(userService)
	router := http.NewServeMux()

	router.HandleFunc("POST /sync", web.HandleWithError(controller.SyncHandler))

	return router
}
