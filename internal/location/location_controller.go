package location

import (
	"iditusi/internal/event"
	"iditusi/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

type LocationController struct {
	path                  string
	addNewLocationUseCase string
	locationStorage       LocationStorage
	eventStorage          event.Storage
}

func NewLocationController(router fiber.Router, locationStorage LocationStorage) *LocationController {
	locationController := &LocationController{
		locationStorage: locationStorage,
	}

	path := router.Group("/locations")
	path.Get("/", locationController.viewLocations)
	path.Post("/", locationController.addNewLocation)
	path.Get("/:id", locationController.viewLocation)
	path.Patch("/:id", locationController.updateLocation)
	return locationController
}

type ViewLocationResponse struct {
	ID          string       `json:"id"`
	Type        LocationType `json:"type"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Stages      []string     `json:"stages"`
	Address     string       `json:"address"`
}

func (c *LocationController) viewLocations(ctx *fiber.Ctx) error {
	return fiber.ErrNotImplemented
}

type AddNewLocationRequest struct {
	Type        string
	Name        string
	Description string
	Stages      []string
	Address     string
}

func (c *LocationController) addNewLocation(ctx *fiber.Ctx) error {
	var request AddNewLocationRequest
	if err := ctx.BodyParser(&request); err != nil {
		return err
	}

	location := NewLocation(utils.NewID(), request.Name, parseLocationType(request.Type), request.Description, request.Address, request.Stages)
	id, err := c.locationStorage.Create(location)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	location.ID = id
	return ctx.JSON(ViewLocationResponse{
		ID:          location.ID,
		Type:        location.Type,
		Name:        location.Name,
		Description: location.Description,
		Stages:      location.Stages,
		Address:     location.Address,
	})
}
func (c *LocationController) viewLocation(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	// id, err := strconv.Atoi(value)
	// if err != nil {
	// 	ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	// }
	location, err := c.locationStorage.Get(id)
	if err != nil {
		ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return ctx.JSON(location)
}
func (c *LocationController) updateLocation(ctx *fiber.Ctx) error {
	return ctx.SendStatus(fiber.StatusNotImplemented)
}
