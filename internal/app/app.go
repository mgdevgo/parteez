package app

import (
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/jackc/pgx/v5/pgxpool"

	"parteez/internal/config"
	"parteez/internal/event"
)

type App struct {
	HTTPServer *fiber.App
	Storage    *pgxpool.Pool
	log        *slog.Logger
	// TokenManager auth.TokenManager
	// Cron cron.Worker
	// Cache redis.Client
}

func New(conf *config.Config, db *pgxpool.Pool, log *slog.Logger) App {
	storage := NewStorage(db, log)

	app := fiber.New(fiber.Config{
		ServerHeader: "Fiber",
		AppName:      "Parteez v0.1.0",
		//DisableStartupMessage: false,
		ReadTimeout:  conf.HTTPServer.Timeout * time.Second,
		WriteTimeout: conf.HTTPServer.Timeout * time.Second,
		IdleTimeout:  conf.HTTPServer.IdleTimeout * time.Second,
	})
	app.Use(limiter.New(), logger.New())

	auth := app.Group("/auth/v1")
	auth.Post("/token", func(ctx *fiber.Ctx) error { return ctx.SendStatus(fiber.StatusNotImplemented) })
	auth.Get("/authorize", func(ctx *fiber.Ctx) error { return ctx.SendStatus(fiber.StatusNotImplemented) })
	auth.Post("/revoke", func(ctx *fiber.Ctx) error { return ctx.SendStatus(fiber.StatusNotImplemented) })

	api := app.Group("/api/v1")

	event.EventRouter(log, api, storage.events)

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

	return App{
		HTTPServer: app,
		Storage:    db,
		log:        log,
	}
}
