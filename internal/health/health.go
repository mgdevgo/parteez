package health

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

type HealthHandler struct {
}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) checkHealth(ctx *fiber.Ctx) error {
	return ctx.JSON(fiber.Map{
		"status": "OK",
		"time":   time.Now(),
	})
}

func (h *HealthHandler) Register(router fiber.Router) {
	router.Get("/", h.checkHealth)
}

func RegisterRoutes(router fiber.Router) {
	healthHandler := NewHealthHandler()
	healthHandler.Register(router)
}
