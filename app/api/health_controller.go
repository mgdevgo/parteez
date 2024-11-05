package api

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

type HealthController struct {
}

func NewHealthController() *HealthController {
	return &HealthController{}
}

func (hc *HealthController) Boot(routes fiber.Router) {
	health := routes.Group("/health")
	health.Get("/", hc.isHealthy)
}

func (hc *HealthController) isHealthy(ctx *fiber.Ctx) error {
	return ctx.JSONP(fiber.Map{
		"status":    "OK",
		"checkTime": time.Now(),
	})
}
