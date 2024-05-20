package search

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type SearchHandler struct {
	searchService SearchService
}

func NewSearchHandler(searchService SearchService) *SearchHandler {
	return &SearchHandler{
		searchService: searchService,
	}
}

func (h *SearchHandler) SearchForResources(ctx *fiber.Ctx) error {
	term := ctx.Query("term")
	types := ctx.Query("types")

	fmt.Println("Term: " + term)
	fmt.Println("Types: " + types)

	return fiber.ErrNotImplemented
}
