package feedback

import "github.com/gofiber/fiber/v2"

type handler struct{}

func NewHandler() *handler {
	return &handler{}
}

func (h *handler) RegisterRoutes(router fiber.Router) {
	ratings := router.Group("/ratings")
	ratings.Put("/events/:id", func(ctx *fiber.Ctx) error { return fiber.ErrNotImplemented })
	ratings.Delete("/events/:id", func(ctx *fiber.Ctx) error { return fiber.ErrNotImplemented })

	router.Get("/feedback", func(ctx *fiber.Ctx) error { return fiber.ErrNotImplemented })
}
