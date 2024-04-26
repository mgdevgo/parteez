package auth

import "github.com/gofiber/fiber/v2"

type handler struct{}

func NewHandler() *handler {
	return &handler{}
}

func (h *handler) RegisterRoutes(router fiber.Router) {
	auth := router.Group("/auth")
	auth.Post("/token", func(ctx *fiber.Ctx) error { return fiber.ErrNotImplemented })
	auth.Get("/authorize", func(ctx *fiber.Ctx) error { return fiber.ErrNotImplemented })
	auth.Post("/revoke", func(ctx *fiber.Ctx) error { return fiber.ErrNotImplemented })
}
