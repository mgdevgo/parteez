package app

import (
	"context"
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/jackc/pgx/v5/pgxpool"

	"parteez/internal/config"
	"parteez/internal/event"
)

const (
	EnvLocal = "local"
	EnvDev   = "dev"
	EnvProd  = "prod"
)

type App struct {
	HTTPServer *fiber.App
	Storage    *pgxpool.Pool
	// TokenManager auth.TokenManager
	// Cron cron.Worker
	// Cache redis.Client
}

func New(ctx context.Context, config *config.Config, log *slog.Logger) (*App, error) {
	const op = "app.New"
	db, _ := pgxpool.New(ctx, config.DatabaseURL)
	if err := db.Ping(ctx); err != nil {
		log.Error("failed to connect database: %v\n", err)
	}

	app := fiber.New(fiber.Config{
		ServerHeader:          "Fiber",
		AppName:               "Parteez v0.1.0",
		DisableStartupMessage: config.AppEnv != EnvProd,
		ReadTimeout:           config.HTTPServer.Timeout * time.Second,
		WriteTimeout:          config.HTTPServer.Timeout * time.Second,
		IdleTimeout:           config.HTTPServer.IdleTimeout * time.Second,
	})
	app.Use(limiter.New(), logger.New())

	auth := app.Group("/auth/v1")
	auth.Post("/token", func(ctx *fiber.Ctx) error { return ctx.SendStatus(fiber.StatusNotImplemented) })
	auth.Get("/authorize", func(ctx *fiber.Ctx) error { return ctx.SendStatus(fiber.StatusNotImplemented) })
	auth.Post("/revoke", func(ctx *fiber.Ctx) error { return ctx.SendStatus(fiber.StatusNotImplemented) })

	api := app.Group("/api/v1")

	events := api.Group("/events")
	events.Get("/", func(ctx *fiber.Ctx) error { return ctx.SendStatus(fiber.StatusNotImplemented) })
	events.Post("/", event.HandleCreateEvent(nil, nil))
	events.Get("/search", func(ctx *fiber.Ctx) error { return ctx.SendStatus(fiber.StatusNotImplemented) })
	events.Get("/:id", func(ctx *fiber.Ctx) error { return ctx.SendStatus(fiber.StatusNotImplemented) })
	events.Patch("/:id", func(ctx *fiber.Ctx) error { return ctx.SendStatus(fiber.StatusNotImplemented) })

	publications := api.Group("/publications")
	publications.Put("/events/:id", func(ctx *fiber.Ctx) error { return ctx.SendStatus(501) })
	publications.Delete("/events/:id", func(ctx *fiber.Ctx) error { return ctx.SendStatus(501) })

	ratings := api.Group("/ratings")
	ratings.Put("/events/:id", func(ctx *fiber.Ctx) error { return ctx.SendStatus(501) })
	ratings.Delete("/events/:id", func(ctx *fiber.Ctx) error { return ctx.SendStatus(501) })

	me := api.Group("/me")
	me.Get("/feed", func(ctx *fiber.Ctx) error { return ctx.SendStatus(501) })
	me.Get("/library", func(ctx *fiber.Ctx) error { return ctx.SendStatus(501) })
	//me/library?ids[albums]=1577502911
	me.Post("/library", func(ctx *fiber.Ctx) error { return ctx.SendStatus(501) })

	locations := api.Group("/locations")
	locations.Get("/", func(ctx *fiber.Ctx) error { return ctx.SendStatus(501) })
	locations.Post("/", func(ctx *fiber.Ctx) error { return ctx.SendStatus(501) })
	locations.Get("/:id", func(ctx *fiber.Ctx) error { return ctx.SendStatus(501) })
	locations.Patch("/:id", func(ctx *fiber.Ctx) error { return ctx.SendStatus(501) })

	api.Get("/feedback", func(ctx *fiber.Ctx) error { return ctx.SendStatus(501) })
	api.Get("/health", func(ctx *fiber.Ctx) error { return ctx.SendStatus(501) })

	return &App{HTTPServer: app, Storage: db}, nil
}
