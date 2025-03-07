package handlers

import (
	"parteez/internal/domain/venue"

	"github.com/gofiber/fiber/v2"
)

type VenueHandler struct {
	venues venue.VenueRepository
}

func NewVenueHandler(venues venue.VenueRepository) *VenueHandler {
	return &VenueHandler{
		venues: venues,
	}
}

func (h *VenueHandler) Register(router fiber.Router) {
	router.Get("/venue/", h.getVenues)
	router.Get("/venue/:id", h.getVenue)
}

func (h *VenueHandler) getVenues(ctx *fiber.Ctx) error {
	result, err := h.venues.FindAll(ctx.Context())
	if err != nil {
		return err
	}

	return ctx.JSON(result)
}

func (h *VenueHandler) getVenue(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id", 0)
	if err != nil {
		return err
	}

	result, err := h.venues.FindById(ctx.Context(), id)
	if err != nil {
		return err
	}

	return ctx.JSON(result)
}
