package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/recover"
	slogfiber "github.com/samber/slog-fiber"

	"parteez/internal/config"
	eventHandler "parteez/internal/domain/events/handler"
	eventStorage "parteez/internal/domain/events/postgres"
	eventService "parteez/internal/domain/events/service"
	venueHandler "parteez/internal/domain/venue/handler"
	venueStorage "parteez/internal/domain/venue/postgres"
	"parteez/internal/errors"
	"parteez/internal/health"
	"parteez/internal/version"
	"parteez/pkg/environment"
	"parteez/pkg/postgres"
)

func main() {
	if err := run(os.Args); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run(args []string) error {
	appenv := environment.Detect(args)
	_ = appenv

	config := config.Load()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	db, err := postgres.New(config.DatabaseURL)
	if err != nil {
		return err
	}

	fiberConfig := fiber.Config{
		AppName:      fmt.Sprintf("parteez v%s", version.VERSION),
		ReadTimeout:  config.HTTPServer.Timeout * time.Second,
		WriteTimeout: config.HTTPServer.Timeout * time.Second,
		IdleTimeout:  config.HTTPServer.IdleTimeout * time.Second,
		ErrorHandler: errors.NewErrorHandler(logger),
	}
	app := fiber.New(fiberConfig)

	rateLimitConfig := limiter.Config{
		LimitReached: func(ctx *fiber.Ctx) error {
			return &errors.Error{
				Status: fiber.StatusTooManyRequests,
				Code:   string(errors.ErrorCodeRateLimit),
			}
		},
	}

	recoverConfig := recover.Config{
		EnableStackTrace: true,
	}

	app.Use(
		cors.New(),
		helmet.New(),
		limiter.New(rateLimitConfig),
		slogfiber.New(logger),
		recover.New(recoverConfig),
	)

	app.Static("/telegram-mini-app", "./web/mini-app/build")

	eventRepository := eventStorage.NewEventStorage(db, nil)
	venueRepository := venueStorage.NewVenueStorage(db, nil)
	eventCrudService := eventService.NewEventCrudService(eventRepository, venueRepository, nil)

	eventHandler := eventHandler.NewEventHandler(eventRepository, eventCrudService)
	venueHandler := venueHandler.NewVenueHandler(venueRepository)
	healthHandler := health.NewHealthHandler()

	app.Route("/health", healthHandler.Register)

	apiv1 := app.Group("/api/v1")
	apiv1.Route("/events", eventHandler.Register)
	apiv1.Route("/venues", venueHandler.Register)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	defer stop()

	go func() {
		if err := app.Listen(":8080"); err != nil {
			logger.Error("Failed to start http server", "error", err)
		}
	}()

	<-ctx.Done()

	if err := app.Shutdown(); err != nil {
		logger.Error("Could not shutdown", "error", err)
	}

	db.Close()

	return nil
}
