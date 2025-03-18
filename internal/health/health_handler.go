package health

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

func NewHealthHandler() func(router fiber.Router) {
	return func(router fiber.Router) {
		router.Get("/", checkHealth)
	}
}

func checkHealth(ctx *fiber.Ctx) error {
	return ctx.JSON(fiber.Map{
		"status": "OK",
		"time":   time.Now(),
	})
}
