package app

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"iditusi/pkg/auth"
	"iditusi/pkg/event"
	"iditusi/pkg/location"
	"iditusi/pkg/shared/api"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	fiberlogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type application struct {
	version    string
	config     *Config
	httpServer *fiber.App
	database   *pgxpool.Pool
	// cache      *cache.Cache
	logger  *zap.SugaredLogger
	wg      sync.WaitGroup
	context context.Context
}

func Run(ctx context.Context, args []string) error {
	app := new(application)
	app.version = "0.1.0"
	app.config = loadConfig()

	var err error

	app.database, err = pgxpool.New(ctx, app.config.DatabaseURL)
	if err != nil {
		return fmt.Errorf("pgxpool.New: %w", err)
	}

	pingctx, cancel := context.WithTimeout(ctx, time.Second*2)
	defer cancel()
	if err := app.database.Ping(pingctx); err != nil {
		return fmt.Errorf("ping: %w", err)
	}

	app.httpServer = fiber.New(fiber.Config{
		AppName: fmt.Sprintf("iditusi v%s", app.version),
		// DisableStartupMessage: true,
		ReadTimeout:  app.config.HTTPServer.Timeout * time.Second,
		WriteTimeout: app.config.HTTPServer.Timeout * time.Second,
		IdleTimeout:  app.config.HTTPServer.IdleTimeout * time.Second,
		ErrorHandler: handelError(),
	})
	app.httpServer.Use(
		fiberlogger.New(),
		limiter.New(limiter.Config{
			LimitReached: func(ctx *fiber.Ctx) error {
				return fiber.ErrTooManyRequests
			},
		}),
	)

	authHandler := auth.NewHandler()

	api := app.httpServer.Group("/api/v1")

	eventStorage := event.NewStorage(app.database)
	locationStorage := location.NewStorage(app.database)

	eventService := event.NewService(eventStorage, locationStorage)

	eventHandler := event.NewHandler(eventService)
	locationHandler := location.NewHandler(locationStorage)

	locationHandler.SetupRoutes(api)
	eventHandler.RegisterRoutes(api)
	authHandler.SetupRoutes(api)

	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	app.wg.Add(1)
	go func() error {
		defer app.wg.Done()

		if err := app.httpServer.Listen(app.config.HTTPServer.Address); err != nil {
			return fmt.Errorf("httpServer.Listen: %w", err)
		}
		return nil
	}()

	app.wg.Add(1)
	go func() error {
		defer app.wg.Done()

		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		if err := app.httpServer.ShutdownWithContext(ctx); err != nil {
			return fmt.Errorf("httpServer.ShutdownWithContext: %s\n", err)
		}

		app.database.Close()

		return nil
	}()

	app.wg.Wait()

	return nil
}

func handelError() func(ctx *fiber.Ctx, err error) error {
	return func(ctx *fiber.Ctx, err error) error {
		var response api.Error

		var fiberError *fiber.Error
		var apiError api.Error

		switch true {
		case errors.As(err, &fiberError):
			response.Status = fiberError.Code
			response.Code = strings.ReplaceAll(strings.ToUpper(fiberError.Message), " ", "_")
		case errors.As(err, &apiError):
			response = apiError
		default:
			response.Status = fiber.StatusInternalServerError
			response.Code = "INTERNAL_SERVER_ERROR"
			response.Detail = err.Error()
		}

		return ctx.Status(response.Status).JSON(response)
	}
}
