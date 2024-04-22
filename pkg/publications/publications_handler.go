package publications

import "github.com/gofiber/fiber/v2"

type handler struct{}

func NewHandler() *handler {
	return &handler{}
}

func (h *handler) SetupRoutes(router fiber.Router) {
	publications := router.Group("/publications")
	publications.Put("/events/:id", func(ctx *fiber.Ctx) error { return fiber.ErrNotImplemented })
	publications.Delete("/events/:id", func(ctx *fiber.Ctx) error { return fiber.ErrNotImplemented })
}
