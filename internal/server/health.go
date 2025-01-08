package server

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

func (s *Server) handleHealth(ctx *fiber.Ctx) error {
	return ctx.JSONP(fiber.Map{
		"status":    "OK",
		"checkTime": time.Now(),
	})
}
