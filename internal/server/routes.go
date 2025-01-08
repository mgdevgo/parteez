package server

import "github.com/gofiber/fiber/v2"

func (s *Server) configureRoutes() {
	router := s.httpServer

	router.Get("/health", s.handleHealth)

	apiv1 := router.Group("/api/v1")

	apiv1.Get("/events", func(ctx *fiber.Ctx) error { return fiber.ErrNotImplemented })
	apiv1.Get("/events/:id", func(ctx *fiber.Ctx) error { return fiber.ErrNotImplemented })

	apiv1.Get("/locations", func(ctx *fiber.Ctx) error { return fiber.ErrNotImplemented })
	apiv1.Get("/locations/:id", func(ctx *fiber.Ctx) error { return fiber.ErrNotImplemented })

	search := apiv1.Group("/search")
	search.Get("/", func(ctx *fiber.Ctx) error { return fiber.ErrNotImplemented })
	search.Get("/hints", func(ctx *fiber.Ctx) error { return fiber.ErrNotImplemented })
	search.Get("/suggestions", func(ctx *fiber.Ctx) error { return fiber.ErrNotImplemented })
	search.Get("/me/library/search", func(ctx *fiber.Ctx) error { return fiber.ErrNotImplemented })

	apiv1.Post("/auth/token", func(ctx *fiber.Ctx) error { return fiber.ErrNotImplemented })
	apiv1.Post("/auth/revoke", func(ctx *fiber.Ctx) error { return fiber.ErrNotImplemented })

	apiv1.Put("/ratings/events/:id", func(ctx *fiber.Ctx) error { return fiber.ErrNotImplemented })
	apiv1.Delete("/ratings/events/:id", func(ctx *fiber.Ctx) error { return fiber.ErrNotImplemented })

	apiv1.Get("/feedback", func(ctx *fiber.Ctx) error { return fiber.ErrNotImplemented })

	apiv1.Get("/me/feed", func(ctx *fiber.Ctx) error { return fiber.ErrNotImplemented })
	apiv1.Get("/me/library", func(ctx *fiber.Ctx) error { return fiber.ErrNotImplemented })
	// me/library?ids[albums]=1577502911
	apiv1.Post("/me/library", func(ctx *fiber.Ctx) error { return fiber.ErrNotImplemented })

	admin := apiv1.Group("/admin")

	admin.Get("/events", func(ctx *fiber.Ctx) error { return fiber.ErrNotImplemented })
	admin.Post("/events", func(ctx *fiber.Ctx) error { return fiber.ErrNotImplemented })
	admin.Get("/events/:id", func(ctx *fiber.Ctx) error { return fiber.ErrNotImplemented })
	admin.Patch("/events/:id", func(ctx *fiber.Ctx) error { return fiber.ErrNotImplemented })

	admin.Get("/locations", func(ctx *fiber.Ctx) error { return fiber.ErrNotImplemented })
	admin.Post("/locations", func(ctx *fiber.Ctx) error { return fiber.ErrNotImplemented })
	admin.Get("/locations/:id", func(ctx *fiber.Ctx) error { return fiber.ErrNotImplemented })
	admin.Patch("/locations/:id", func(ctx *fiber.Ctx) error { return fiber.ErrNotImplemented })

	admin.Get("/users", func(ctx *fiber.Ctx) error { return fiber.ErrNotImplemented })
	admin.Post("/users", func(ctx *fiber.Ctx) error { return fiber.ErrNotImplemented })
	admin.Get("/users/:id", func(ctx *fiber.Ctx) error { return fiber.ErrNotImplemented })
	admin.Patch("/users/:id", func(ctx *fiber.Ctx) error { return fiber.ErrNotImplemented })

	// ?ids[events]=123,1234&ids[locations]=123,1234
	admin.Post("/publications", func(ctx *fiber.Ctx) error { return fiber.ErrNotImplemented })
	admin.Delete("/publications/:id", func(ctx *fiber.Ctx) error { return fiber.ErrNotImplemented })
}
