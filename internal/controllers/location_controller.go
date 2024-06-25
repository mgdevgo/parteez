package controllers

import (
	"iditusi/internal/repositories"

	"github.com/gofiber/fiber/v2"
)

type LocationController struct {
	eventRepository    repositories.EventRepository
	locationRepository repositories.VenueRepository
}

func NewLocationController(eventRepository repositories.EventRepository) *LocationController {
	return &LocationController{
		eventRepository: eventRepository,
	}
}

func (lc *LocationController) Boot(routes fiber.Router) {
	locations := routes.Group("/locations")
	locations.Get("/", lc.getLocations)
	locations.Get("/:id", lc.getLocation)
}

func (lc *LocationController) getLocations(ctx *fiber.Ctx) error {
	result, err := lc.locationRepository.FindAll(ctx.Context())
	if err != nil {
		return err
	}

	return ctx.JSON(result)
}

func (lc *LocationController) getLocation(ctx *fiber.Ctx) error {
	id, _ := ctx.ParamsInt("id", 0)
	result, err := lc.locationRepository.FindById(ctx.Context(), id)
	if err != nil {
		return err
	}

	return ctx.JSON(result)
}
