package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"iditusi/pkg/api"
	"iditusi/pkg/events"
	"iditusi/pkg/health"
	"iditusi/pkg/listing"
	"iditusi/pkg/search"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/jackc/pgx/v5/pgxpool"
)

const version = "0.1.1"

func Run(ctx context.Context, args []string) error {
	config := loadConfig()

	var err error

	db, err := pgxpool.New(ctx, config.DatabaseURL)
	if err != nil {
		return fmt.Errorf("pgxpool.New: %w", err)
	}

	pingctx, cancel := context.WithTimeout(ctx, time.Second*2)
	defer cancel()
	if err := db.Ping(pingctx); err != nil {
		return fmt.Errorf("ping: %w", err)
	}

	// Storage
	eventStorage := events.NewPostgresEventStorage(db)
	// locationsStorage := locations.NewPostgresLocationsStorage(db)

	// Services
	searchService := search.NewSearchService()
	eventListingService := listing.NewEventListingService(eventStorage)
	// locationListingService := listing.NewLocationListingService(locationsStorage)
	// userManagementService := administration.NewUserManagementService()
	// eventManagementService := administration.NewEventManagementService(eventsStorage)

	// Handlers
	searchHandler := search.NewSearchHandler(searchService)
	healthHandler := health.NewHealthHandler()
	listingHandler := listing.NewListingHandler(eventListingService)

	app := fiber.New(fiber.Config{
		AppName: fmt.Sprintf("iditusi v%s", version),
		// DisableStartupMessage: true,
		ReadTimeout:  config.HTTPServer.Timeout * time.Second,
		WriteTimeout: config.HTTPServer.Timeout * time.Second,
		IdleTimeout:  config.HTTPServer.IdleTimeout * time.Second,
		ErrorHandler: api.ErrorHandler(),
	})
	app.Use(
		logger.New(),
		limiter.New(limiter.Config{
			LimitReached: func(ctx *fiber.Ctx) error {
				return &api.Error{
					Status: fiber.StatusTooManyRequests,
					Code:   string(api.ErrorCodeRateLimit),
				}
				// return fiber.ErrTooManyRequests
			},
		}),
	)

	apiv1 := app.Group("/api/v1")
	admin := apiv1.Group("/admin")

	app.Get("/health", healthHandler.CheckHealth)

	app.Get("/auth", func(ctx *fiber.Ctx) error { return fiber.ErrNotImplemented })
	app.Post("/auth/token", func(ctx *fiber.Ctx) error { return fiber.ErrNotImplemented })
	app.Post("/auth/revoke", func(ctx *fiber.Ctx) error { return fiber.ErrNotImplemented })

	apiv1.Get("/events", listingHandler.GetMultipleEvents)
	apiv1.Get("/events/:id", listingHandler.GetEventByID)

	apiv1.Get("/locations", listingHandler.GetAllLocations)
	apiv1.Get("/locations/:id", listingHandler.GetLocationByID)

	apiv1.Get("/search", searchHandler.SearchForResources)
	apiv1.Get("/search/hints", func(ctx *fiber.Ctx) error { return fiber.ErrNotImplemented })
	apiv1.Get("/search/suggestions", func(ctx *fiber.Ctx) error { return fiber.ErrNotImplemented })
	apiv1.Get("/me/library/search", func(ctx *fiber.Ctx) error { return fiber.ErrNotImplemented })

	apiv1.Put("/ratings/events/:id", func(ctx *fiber.Ctx) error { return fiber.ErrNotImplemented })
	apiv1.Delete("/ratings/events/:id", func(ctx *fiber.Ctx) error { return fiber.ErrNotImplemented })

	apiv1.Get("/feedback", func(ctx *fiber.Ctx) error { return fiber.ErrNotImplemented })

	apiv1.Get("/me/feed", func(ctx *fiber.Ctx) error { return fiber.ErrNotImplemented })
	apiv1.Get("/me/library", func(ctx *fiber.Ctx) error { return fiber.ErrNotImplemented })
	// me/library?ids[albums]=1577502911
	apiv1.Post("/me/library", func(ctx *fiber.Ctx) error { return fiber.ErrNotImplemented })

	admin.Get("/events", func(ctx *fiber.Ctx) error { return fiber.ErrNotImplemented })
	admin.Post("/events", func(ctx *fiber.Ctx) error { return fiber.ErrNotImplemented })
	admin.Get("/events/:id", func(ctx *fiber.Ctx) error { return fiber.ErrNotImplemented })
	admin.Patch("/events/:id", func(ctx *fiber.Ctx) error { return fiber.ErrNotImplemented })

	admin.Get("/locations", func(ctx *fiber.Ctx) error { return fiber.ErrNotImplemented })
	admin.Post("/locations", func(ctx *fiber.Ctx) error { return fiber.ErrNotImplemented })
	admin.Get("/locations/:id", func(ctx *fiber.Ctx) error { return fiber.ErrNotImplemented })
	admin.Patch("/locations/:id", func(ctx *fiber.Ctx) error { return fiber.ErrNotImplemented })

	admin.Get("/users", func(ctx *fiber.Ctx) error { return fiber.ErrNotImplemented })
	admin.Post("/users", func(ctx *fiber.Ctx) error { return fiber.ErrNotImplemented })
	admin.Get("/users/:id", func(ctx *fiber.Ctx) error { return fiber.ErrNotImplemented })
	admin.Patch("/users/:id", func(ctx *fiber.Ctx) error { return fiber.ErrNotImplemented })

	// ?ids[events]=123,1234&ids[locations]=123,1234
	admin.Post("/publications", func(ctx *fiber.Ctx) error { return fiber.ErrNotImplemented })
	admin.Delete("/publications/:id", func(ctx *fiber.Ctx) error { return fiber.ErrNotImplemented })

	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() error {
		defer wg.Done()

		if err := app.Listen(config.HTTPServer.Address); err != nil {
			return fmt.Errorf("httpServer.Listen: %w", err)
		}
		return nil
	}()

	wg.Add(1)
	go func() error {
		defer wg.Done()

		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		if err := app.ShutdownWithContext(ctx); err != nil {
			return fmt.Errorf("httpServer.ShutdownWithContext: %s\n", err)
		}

		db.Close()

		return nil
	}()

	wg.Wait()

	return nil
}
