package api

import (
	"fmt"

	"iditusi/internal/services"

	"github.com/gofiber/fiber/v2"
)

type SearchController struct {
	searchService services.SearchService
}

func NewSearchController(searchService services.SearchService) *SearchController {
	return &SearchController{
		searchService: searchService,
	}
}

func (sc *SearchController) Boot(routes fiber.Router) {
	search := routes.Group("/search")
	search.Get("/", sc.search)
	search.Get("/hints", sc.hints)
	search.Get("/suggestions", sc.suggestions)

	library := search.Group("/me/library")
	library.Get("/me/library/search", sc.librarySearch)
}

func (sc *SearchController) search(ctx *fiber.Ctx) error {
	term := ctx.Query("term")
	types := ctx.Query("types")

	fmt.Println("Term: " + term)
	fmt.Println("Types: " + types)

	return fiber.ErrNotImplemented
}

func (sc *SearchController) hints(ctx *fiber.Ctx) error {
	return fiber.ErrNotImplemented
}

func (sc *SearchController) suggestions(ctx *fiber.Ctx) error {
	return fiber.ErrNotImplemented
}

func (sc *SearchController) librarySearch(ctx *fiber.Ctx) error {
	return fiber.ErrNotImplemented
}
