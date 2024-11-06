package main

import (
	"context"
	"fmt"
	"iditusi/app/api"
	"iditusi/app/api/response"
	"iditusi/app/config"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/cobra"
)

func NewServeCommand() *cobra.Command {
	serveCommand := &cobra.Command{
		Use: "serve",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			config := config.Load()

			var err error

			db, err := pgxpool.New(ctx, config.DatabaseURL)
			if err != nil {
				return fmt.Errorf("pgxpool.New: %w", err)
			}

			pingContext, cancelPing := context.WithTimeout(ctx, time.Second*2)
			defer cancelPing()
			if err := db.Ping(pingContext); err != nil {
				return fmt.Errorf("ping: %w", err)
			}

			app := fiber.New(fiber.Config{
				AppName: fmt.Sprintf("iditusi v%s", version),
				// DisableStartupMessage: true,
				ReadTimeout:  config.HTTPServer.Timeout * time.Second,
				WriteTimeout: config.HTTPServer.Timeout * time.Second,
				IdleTimeout:  config.HTTPServer.IdleTimeout * time.Second,
				ErrorHandler: api.ErrorHandler,
			})

			app.Use(logger.New())

			limiterMiddlewareConfig := limiter.Config{
				LimitReached: func(ctx *fiber.Ctx) error {
					return &response.Error{
						Status: fiber.StatusTooManyRequests,
						Code:   string(response.ErrorCodeRateLimit),
					}
					// return fiber.ErrTooManyRequests
				},
			}
			app.Use(limiter.New(limiterMiddlewareConfig))

			// eventRepository := postgres.NewPostgresEventRepository(db)
			// locationRepository := postgres.NewLocationRepository(db)

			apiv1 := app.Group("/api/v1")
			// controllers.NewHealthController().Boot(apiv1)
			// controllers.NewEventController(eventRepository).Boot(apiv1)
			// controllers.NewLocationController(eventRepository).Boot(apiv1)
			// controllers.NewSearchController(eventRepository).Boot(apiv1)

			app.Get("/auth", func(ctx *fiber.Ctx) error { return fiber.ErrNotImplemented })
			app.Post("/auth/token", func(ctx *fiber.Ctx) error { return fiber.ErrNotImplemented })
			app.Post("/auth/revoke", func(ctx *fiber.Ctx) error { return fiber.ErrNotImplemented })

			apiv1.Put("/ratings/events/:id", func(ctx *fiber.Ctx) error { return fiber.ErrNotImplemented })
			apiv1.Delete("/ratings/events/:id", func(ctx *fiber.Ctx) error { return fiber.ErrNotImplemented })

			apiv1.Get("/feedback", func(ctx *fiber.Ctx) error { return fiber.ErrNotImplemented })

			apiv1.Get("/me/feed", func(ctx *fiber.Ctx) error { return fiber.ErrNotImplemented })
			apiv1.Get("/me/library", func(ctx *fiber.Ctx) error { return fiber.ErrNotImplemented })
			// me/library?ids[albums]=1577502911
			apiv1.Post("/me/library", func(ctx *fiber.Ctx) error { return fiber.ErrNotImplemented })

			admin := apiv1.Group("/admin")

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
		},
	}

	serveCommand.Flags().String("bind", "0.0.0.0:8080", "Convenience for setting hostname and port together.")
	serveCommand.Flags().StringP("hostname", "h", "0.0.0.0", "Set the hostname the server will run on.")
	serveCommand.Flags().StringP("port", "p", "8080", "Set the port the server will run on.")

	return serveCommand
}
