package health

import "github.com/gofiber/fiber/v2"

type handler struct{}

func NewHandler() *handler {
	return &handler{}
}

func (h *handler) SetupRoutes(router fiber.Router) {
	router.Get("/health", func(ctx *fiber.Ctx) error { return fiber.ErrNotImplemented })

}
