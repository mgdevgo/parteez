package handlers

import (
	"iditusi/internal/location"
	"log/slog"

	"github.com/gofiber/fiber/v2"
)

func Register(router fiber.Router, log *slog.Logger, storage location.Storage) {
	locations := router.Group("/locations")
	locations.Get("/", func(ctx *fiber.Ctx) error { return ctx.SendStatus(501) })
	locations.Post("/", func(ctx *fiber.Ctx) error { return ctx.SendStatus(501) })
	locations.Get("/:id", func(ctx *fiber.Ctx) error { return ctx.SendStatus(501) })
	locations.Patch("/:id", func(ctx *fiber.Ctx) error { return ctx.SendStatus(501) })
}
