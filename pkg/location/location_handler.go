package location

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type handler struct {
	locationStorage Storage
}

func NewHandler(locationStorage Storage) *handler {
	return &handler{
		locationStorage: locationStorage,
	}
}

func (h *handler) RegisterRoutes(router fiber.Router) {
	locations := router.Group("/locations")
	locations.Get("/", h.viewLocations)
	locations.Post("/", h.addNewLocation)
	locations.Get("/:id", h.viewLocation)
	locations.Patch("/:id", h.updateLocation)
}

type viewLocationResponse struct {
	ID          int      `json:"id"`
	Type        Kind     `json:"type"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Stages      []string `json:"stages"`
	Address     string   `json:"address"`
}

func (h *handler) viewLocations(ctx *fiber.Ctx) error {
	locations, err := h.locationStorage.FindAll()
	if err != nil {
		return err
	}
	return ctx.JSON(locations)
}

func (h *handler) addNewLocation(ctx *fiber.Ctx) error {
	type addNewLocationRequest struct {
		Type        string
		Name        string
		Description string
		Stages      []string
		Address     string
	}

	request := addNewLocationRequest{}
	if err := ctx.BodyParser(&request); err != nil {
		return err
	}

	location := NewLocation(request.Name)
	location, err := h.locationStorage.Save(location)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(viewLocationResponse{
		ID:          location.ID,
		Type:        location.Type,
		Name:        location.Name,
		Description: location.Description,
		Stages:      location.Stages,
		Address:     location.Address,
	})
}
func (h *handler) viewLocation(ctx *fiber.Ctx) error {
	paramID := ctx.Params("id")
	id, err := strconv.Atoi(paramID)
	if err != nil {
		ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	location, err := h.locationStorage.FindByID(id)
	if err != nil {
		ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return ctx.JSON(location)
}
func (h *handler) updateLocation(ctx *fiber.Ctx) error {
	return ctx.SendStatus(fiber.StatusNotImplemented)
}
