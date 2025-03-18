package handlers

import (
	"parteez/internal/domain/venue"

	"github.com/gofiber/fiber/v2"
)

func NewVenueHandler(venues venue.VenueRepository) func(router fiber.Router) {
	return func(router fiber.Router) {
		router.Get("/venue/", handleGetVenues(venues))
		router.Get("/venue/:id", handleGetVenue(venues))
	}
}

func handleGetVenues(repository venue.VenueRepository) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		result, err := repository.FindAll(ctx.Context())
		if err != nil {
			return err
		}
		return ctx.JSON(result)
	}
}

func handleGetVenue(repository venue.VenueRepository) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		id, err := ctx.ParamsInt("id", 0)
		if err != nil {
			return err
		}

		result, err := repository.FindById(ctx.Context(), id)
		if err != nil {
			return err
		}

		return ctx.JSON(result)
	}
}
