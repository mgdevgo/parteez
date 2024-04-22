package user

import "github.com/gofiber/fiber/v2"

type handler struct{}

func NewHandler() *handler {
	return &handler{}
}

func (h *handler) SetupRoutes(router fiber.Router) {
	me := router.Group("/me")
	me.Get("/feed", func(ctx *fiber.Ctx) error { return fiber.ErrNotImplemented })
	me.Get("/library", func(ctx *fiber.Ctx) error { return fiber.ErrNotImplemented })
	// me/library?ids[albums]=1577502911
	me.Post("/library", func(ctx *fiber.Ctx) error { return fiber.ErrNotImplemented })

	users := router.Group("/users")
	users.Get("/", func(ctx *fiber.Ctx) error { return fiber.ErrNotImplemented })
	users.Post("/", func(ctx *fiber.Ctx) error { return fiber.ErrNotImplemented })
	users.Get("/:id", func(ctx *fiber.Ctx) error { return fiber.ErrNotImplemented })
	users.Patch("/:id", func(ctx *fiber.Ctx) error { return fiber.ErrNotImplemented })
}
