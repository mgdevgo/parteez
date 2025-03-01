package handlers

import (
	"parteez/internal/domain/venue"

	"github.com/gofiber/fiber/v2"
)

type VenueHandler struct {
	venues venue.VenueRepository
}

func NewVenueHandler() *VenueHandler {
	return &VenueHandler{}
}

func (h *VenueHandler) Register(router fiber.Router) {
	router.Get("/", h.getLocationCollection)
	router.Get("/:id", h.getLocationById)
}

func (h *VenueHandler) getLocationCollection(ctx *fiber.Ctx) error {
	result, err := h.venues.FindAll(ctx.Context())
	if err != nil {
		return err
	}

	return ctx.JSON(result)
}

func (h *VenueHandler) getLocationById(ctx *fiber.Ctx) error {
	id, _ := ctx.ParamsInt("id", 0)
	result, err := h.venues.FindById(ctx.Context(), id)
	if err != nil {
		return err
	}

	return ctx.JSON(result)
}
