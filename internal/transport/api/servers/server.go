package servers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/kr0uch/studspace/internal/core/service"
	"github.com/kr0uch/studspace/internal/transport/api/middlewares"
	"github.com/kr0uch/studspace/internal/transport/api/routers"
	"github.com/kr0uch/studspace/pkg/logger"
	"github.com/kr0uch/studspace/pkg/web"
)

type ServerConfig struct {
	Host   string `yaml:"host" env:"API_HOST" env-default:"0.0.0.0"`
	Port   int    `yaml:"port" env:"API_PORT" env-default:"8080"`
	Prefix string `yaml:"prefix" env:"API_PREFIX" env-default:"/api/v1"`
}

type VKConfig struct {
	AppSecret string `yaml:"app_secret" env:"VK_APP_SECRET" env-required:"true"`
}

type Server struct {
	config ServerConfig
	server *http.Server
	logger logger.Logger
}

func NewServer(
	config ServerConfig,
	vkConfig VKConfig,
	logger logger.Logger,
	userService *service.UserService,
	spaceService *service.SpaceService,
	bookingService *service.BookingService,
	eventService *service.EventService,
) *Server {
	vk := middlewares.NewVKAuth(vkConfig.AppSecret, userService)

	authRouter := routers.NewAuthRouter(userService)
	userRouter := routers.NewUserRouter(userService)
	homeRouter := routers.NewHomeRouter(bookingService, spaceService)
	spaceRouter := routers.NewSpaceRouter(spaceService)
	scheduleRouter := routers.NewScheduleRouter(spaceService, bookingService, eventService)
	bookingRouter := routers.NewBookingRouter(bookingService, spaceService, userService)
	eventRouter := routers.NewEventRouter(eventService, spaceService)
	adminRouter := routers.NewAdminRouter(bookingService, eventService, spaceService, userService)

	protected := vk.LoadUser
	apiMux := http.NewServeMux()
	apiMux.Handle("/auth/", http.StripPrefix("/auth", authRouter))
	apiMux.Handle("/me", protected(userRouter))
	apiMux.Handle("/me/", protected(userRouter))
	apiMux.Handle("/home", protected(homeRouter))
	apiMux.Handle("/spaces", protected(spaceRouter))
	apiMux.Handle("/spaces/", protected(spaceRouter))
	apiMux.Handle("/schedule", protected(scheduleRouter))
	apiMux.Handle("/calendar/", protected(scheduleRouter))
	apiMux.Handle("/bookings", protected(bookingRouter))
	apiMux.Handle("/bookings/", protected(bookingRouter))
	apiMux.Handle("/events", protected(eventRouter))
	apiMux.Handle("/events/", protected(eventRouter))
	apiMux.Handle("/admin/", protected(adminRouter))

	handler := web.GlobalMiddleware(
		web.ErrorMiddleware(
			web.LoggerMiddleware(logger)(
				vk.VerifySign(apiMux),
			),
		),
	)

	mainMux := http.NewServeMux()
	mainMux.Handle(config.Prefix+"/", http.StripPrefix(config.Prefix, handler))

	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", config.Host, config.Port),
		Handler: mainMux,
	}

	return &Server{
		config: config,
		server: server,
		logger: logger,
	}
}

func (s *Server) Start() error {
	s.logger.With(logger.Time("started_at", time.Now())).Info(fmt.Sprintf("Server started on %s", s.server.Addr))
	return s.server.ListenAndServe()
}

func (s *Server) Stop() error {
	s.logger.With(logger.Time("stopped_at", time.Now())).Info("Server stopped")
	return s.server.Shutdown(context.Background())
}
