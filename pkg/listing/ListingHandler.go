package listing

import (
	"context"

	"iditusi/pkg/core"

	"github.com/gofiber/fiber/v2"
)

type EventListingService interface {
	ListEvents(ctx context.Context, fromDate string, toDate string) ([]core.Event, error)
	ListEventById(ctx context.Context, id int) (core.Event, error)
}
type LocationListingService interface {
	ListLocations(ctx context.Context) ([]string, error)
	ListLocationById(ctx context.Context, id int) (string, error)
}

type ListingHandler struct {
	eventService    EventListingService
	locationService LocationListingService
}

func NewListingHandler(eventService EventListingService) *ListingHandler {
	return &ListingHandler{
		eventService: eventService,
	}
}

func (h *ListingHandler) GetMultipleEvents(ctx *fiber.Ctx) error {
	// limit := ctx.QueryInt("limit", 5)
	// offset := ctx.QueryInt("offset", 0)
	fromDate := ctx.Query("fromDate", "")
	toDate := ctx.Query("toDate", "")

	result, err := h.eventService.ListEvents(ctx.Context(), fromDate, toDate)
	if err != nil {
		return err
	}

	return ctx.JSON(result)
}

func (h *ListingHandler) GetEventByID(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id", 0)
	if err != nil {
		return err
	}

	result, err := h.eventService.ListEventById(ctx.Context(), id)
	if err != nil {
		return err
	}

	return ctx.JSON(result)
}

func (h *ListingHandler) GetAllLocations(ctx *fiber.Ctx) error {
	result, err := h.locationService.ListLocations(ctx.Context())
	if err != nil {
		return err
	}

	return ctx.JSON(result)
}

func (h *ListingHandler) GetLocationByID(ctx *fiber.Ctx) error {
	id, _ := ctx.ParamsInt("id", 0)
	result, err := h.locationService.ListLocationById(ctx.Context(), id)
	if err != nil {
		return err
	}

	return ctx.JSON(result)
}
