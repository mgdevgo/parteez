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

func (h *HealthHandler) CheckHealth(ctx *fiber.Ctx) error {
	return ctx.JSONP(fiber.Map{
		"status":    "OK",
		"checkTime": time.Now(),
	})
}
