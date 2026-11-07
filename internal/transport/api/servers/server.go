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

type SwaggerConfig struct {
	Enabled bool   `yaml:"enabled" env:"SWAGGER_ENABLED" env-default:"false"`
	Host    string `yaml:"host" env:"SWAGGER_HOST" env-default:"localhost:8080"`
	Schema  string `yaml:"schema" env:"SWAGGER_SCHEMA" env-default:"http"`
}

type ServerConfig struct {
	Host          string        `yaml:"host" env:"API_HOST" env-default:"0.0.0.0"`
	Port          int           `yaml:"port" env:"API_PORT" env-default:"8080"`
	Prefix        string        `yaml:"prefix" env:"API_PREFIX" env-default:"/api/v1"`
	CORSOrigins   []string      `yaml:"cors_origins" env:"CORS_ORIGINS" env-separator:","`
	SwaggerConfig SwaggerConfig `yaml:"swagger"`
}

type VKConfig struct {
	AppSecret   string `yaml:"app_secret" env:"VK_APP_SECRET" env-required:"true"`
	NotifyToken string `yaml:"notify_token" env:"VK_NOTIFY_TOKEN"`
}

const (
	readHeaderTimeout = 5 * time.Second
	readTimeout       = 10 * time.Second
	writeTimeout      = 30 * time.Second
	idleTimeout       = 120 * time.Second
	shutdownTimeout   = 30 * time.Second
)

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
) *Server {
	vk := middlewares.NewVKAuth(vkConfig.AppSecret, userService)

	authRouter := routers.NewAuthRouter(userService)
	userRouter := routers.NewUserRouter(userService)
	homeRouter := routers.NewHomeRouter(bookingService, spaceService)
	spaceRouter := routers.NewSpaceRouter(spaceService)
	scheduleRouter := routers.NewScheduleRouter(spaceService, bookingService)
	bookingRouter := routers.NewBookingRouter(bookingService, spaceService, userService)
	adminRouter := routers.NewAdminRouter(bookingService, spaceService, userService)

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
	apiMux.Handle("/admin/", protected(adminRouter))

	handler := web.CORS(config.CORSOrigins)(
		web.GlobalMiddleware(
			web.ErrorMiddleware(
				web.LoggerMiddleware(logger)(
					vk.VerifySign(apiMux),
				),
			),
		),
	)

	mainMux := http.NewServeMux()
	mainMux.HandleFunc(config.Prefix+"/health", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	if config.SwaggerConfig.Enabled {
		mainMux.Handle(config.Prefix+"/swagger/", middlewares.Swagger(config.Prefix))
	} else {
		mainMux.HandleFunc(config.Prefix+"/swagger/", func(w http.ResponseWriter, r *http.Request) {
			http.NotFound(w, r)
		})
	}
	mainMux.Handle(config.Prefix+"/", http.StripPrefix(config.Prefix, handler))

	server := &http.Server{
		Addr:              fmt.Sprintf("%s:%d", config.Host, config.Port),
		Handler:           mainMux,
		ReadHeaderTimeout: readHeaderTimeout,
		ReadTimeout:       readTimeout,
		WriteTimeout:      writeTimeout,
		IdleTimeout:       idleTimeout,
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
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	return s.server.Shutdown(ctx)
}
