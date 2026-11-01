package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kr0uch/studspace/docs"
	"github.com/kr0uch/studspace/internal/config"
	"github.com/kr0uch/studspace/internal/core/repository/postgres"
	"github.com/kr0uch/studspace/internal/core/service"
	"github.com/kr0uch/studspace/internal/transport/api/servers"
	dbpg "github.com/kr0uch/studspace/pkg/database/postgres"
	"github.com/kr0uch/studspace/pkg/logger"
	"github.com/kr0uch/studspace/pkg/vk"
)

// @title StudSpace API
// @version 1.0
// @description Бэкенд VK Mini App: бронь студенческих залов и запись на мероприятия.
// @host localhost:8080
// @schemes http
// @BasePath /api/v1
//
// @securityDefinitions.apikey VKAuth
// @in header
// @name Authorization
// @description Подпись launch params. Формат: VK vk_user_id=...&vk_ts=...&sign=...
func main() {
	ctx := context.Background()

	cfg, err := config.NewConfig()
	if err != nil {
		panic(err)
	}

	docs.SwaggerInfo.Host = cfg.ServerConfig.SwaggerConfig.Host
	docs.SwaggerInfo.Schemes = []string{cfg.ServerConfig.SwaggerConfig.Schema}
	docs.SwaggerInfo.BasePath = cfg.ServerConfig.Prefix

	zapLogger, err := logger.NewLogger(cfg.LoggerConfig)
	if err != nil {
		panic(err)
	}

	postgresDB, err := dbpg.NewPostgresDB(ctx, cfg.PostgresConfig)
	if err != nil {
		zapLogger.Error("postgres.NewPostgresDB", logger.Error(err))
		return
	}

	userRepo := postgres.NewUserRepository(postgresDB)
	spaceRepo := postgres.NewSpaceRepository(postgresDB)
	bookingRepo := postgres.NewBookingRepository(postgresDB)
	eventRepo := postgres.NewEventRepository(postgresDB)

	notifyRepo := postgres.NewNotifyRepository(vk.NewVKClient(cfg.VK.NotifyToken))

	userService := service.NewUserService(userRepo)
	spaceService := service.NewSpaceService(spaceRepo, bookingRepo, eventRepo)
	bookingService := service.NewBookingService(bookingRepo, eventRepo, spaceRepo, userRepo, notifyRepo)
	eventService := service.NewEventService(eventRepo, bookingRepo, spaceRepo, userRepo, notifyRepo)

	server := servers.NewServer(
		cfg.ServerConfig,
		cfg.VK,
		zapLogger,
		userService,
		spaceService,
		bookingService,
		eventService,
	)

	graceChan := make(chan os.Signal, 1)
	signal.Notify(graceChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := server.Start(); err != nil {
			zapLogger.Error("server.Start", logger.Error(err))
		}
	}()
	<-graceChan

	if err = server.Stop(); err != nil {
		zapLogger.Error("server.Stop", logger.Error(err))
	}

	if err = postgresDB.Close(); err != nil {
		zapLogger.Error("postgresDB.Close", logger.Error(err))
	} else {
		zapLogger.With(logger.Time("stopped_at", time.Now())).Info("postgresDB closed")
	}
}
