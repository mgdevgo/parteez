package application

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"iditusi/internal/location"
	"iditusi/internal/shared/api"

	requestLogger "github.com/gofiber/fiber/v2/middleware/logger"
	"golang.org/x/sync/errgroup"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Run() {
	rootctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	config := newConfig()
	logger := newLogger(config.AppEnv)

	db, err := pgxpool.New(rootctx, config.DatabaseURL)
	if err != nil {
		logger.Error(fmt.Sprintf("pgxpool.New: %w", err))
	}

	pingctx, cancel := context.WithTimeout(rootctx, time.Second*2)
	defer cancel()
	if err := db.Ping(pingctx); err != nil {
		panic(fmt.Errorf("ping: %w", err))
	}

	logger.Info("Connected to database", slog.String("connection_url", config.DatabaseURL))

	server := fiber.New(fiber.Config{
		ServerHeader: "Fiber",
		AppName:      "iditusi v0.1.0",
		// DisableStartupMessage: false,
		ReadTimeout:  config.HTTPServer.Timeout * time.Second,
		WriteTimeout: config.HTTPServer.Timeout * time.Second,
		IdleTimeout:  config.HTTPServer.IdleTimeout * time.Second,
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			response := api.ErrorsResponse{
				Errors: []api.Error{
					{
						Status: fiber.StatusInternalServerError,
						Code:   "INTERNAL_SERVER_ERROR",
					},
				},
			}

			// Retrieve the custom status code if it's a *fiber.Error
			var e *fiber.Error
			if errors.As(err, &e) {
				response.Errors[0].Code = strings.ReplaceAll(strings.ToUpper(e.Message), " ", "_")
				response.Errors[0].Status = e.Code
			}

			return ctx.Status(response.Errors[0].Status).JSON(response)

			// Return from handler
			return nil
		},
	})
	server.Use(
		requestLogger.New(),
		limiter.New(limiter.Config{
			LimitReached: func(ctx *fiber.Ctx) error {
				return fiber.ErrTooManyRequests
			},
		}),
	)

	auth := server.Group("/auth/v1")
	auth.Post("/token", func(ctx *fiber.Ctx) error { return fiber.ErrNotImplemented })
	auth.Get("/authorize", func(ctx *fiber.Ctx) error { return ctx.SendStatus(fiber.StatusNotImplemented) })
	auth.Post("/revoke", func(ctx *fiber.Ctx) error { return ctx.SendStatus(fiber.StatusNotImplemented) })

	api := server.Group("/api/v1")

	// storage := defaultStorage(db, logger)
	// handlers.EventRouter(logger, api, storage.events)

	locationStorage := location.NewLocationStorage(db)
	location.NewLocationController(api, locationStorage)

	publications := api.Group("/publications")
	publications.Put("/events/:id", func(ctx *fiber.Ctx) error { return ctx.SendStatus(501) })
	publications.Delete("/events/:id", func(ctx *fiber.Ctx) error { return ctx.SendStatus(501) })

	ratings := api.Group("/ratings")
	ratings.Put("/events/:id", func(ctx *fiber.Ctx) error { return ctx.SendStatus(501) })
	ratings.Delete("/events/:id", func(ctx *fiber.Ctx) error { return ctx.SendStatus(501) })

	me := api.Group("/me")
	me.Get("/feed", func(ctx *fiber.Ctx) error { return ctx.SendStatus(501) })
	me.Get("/library", func(ctx *fiber.Ctx) error { return ctx.SendStatus(501) })
	// me/library?ids[albums]=1577502911
	me.Post("/library", func(ctx *fiber.Ctx) error { return ctx.SendStatus(501) })

	api.Get("/feedback", func(ctx *fiber.Ctx) error { return ctx.SendStatus(501) })
	api.Get("/health", func(ctx *fiber.Ctx) error { return ctx.SendStatus(501) })

	g, run := errgroup.WithContext(rootctx)
	g.Go(func() error {
		if err := server.Listen(config.HTTPServer.Address); err != nil {
			logger.Error("HTTP server listen error", slog.String("error", err.Error()))
			return err
		}
		return nil
	})

	g.Go(func() error {
		pid := syscall.Getpid()
		for {
			select {
			case <-run.Done():
				logger.Info("Background context for run operation closed - Shutting down...", slog.Int("pid", pid), slog.String("error", run.Err().Error()))
				ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)

				if err := server.ShutdownWithContext(ctx); err != nil {
					logger.Error("Error when shutdown HTTP server", slog.String("error", err.Error()))
				}
				logger.Info("HTTP server stopped")
				cancel()

				db.Close()
				logger.Info("Database connections closed")
				return nil
			}

		}
	})

	// g.Go(func() error {
	//	<-run.Done()
	//
	//	log.Info("Database connections closed")
	//	return nil
	// })
	g.Wait()
	// if err := g.Wait(); err != nil {
	//	log.Error("Application stopped with error", slog.String("error", err.Error()))
	// }

	logger.Info("Application stopped")

}
