package handlers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type SearchHandler struct {
}

func NewSearchController() *SearchHandler {
	return &SearchHandler{}
}

func (h *SearchHandler) Register(router fiber.Router) {
	// place - hints, suggestions, library
	router.Get("/:place", h.search)
}

func (h *SearchHandler) search(ctx *fiber.Ctx) error {
	term := ctx.Query("term")
	types := ctx.Query("types")

	fmt.Println("Term: " + term)
	fmt.Println("Types: " + types)

	return fiber.ErrNotImplemented
}
