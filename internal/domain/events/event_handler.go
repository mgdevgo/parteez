package events

import (
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"

	"parteez/internal/domain/shared"
)

type EventHandler struct {
	eventRepository  EventRepository
	eventCrudService EventCrudService
}

func NewEventHandler(eventRepository EventRepository, eventCrudService EventCrudService) *EventHandler {
	return &EventHandler{
		eventRepository:  eventRepository,
		eventCrudService: eventCrudService,
	}
}

func (h *EventHandler) Register(router fiber.Router) {
	router.Get("/", h.handleGetManyEvents)
	router.Get("/:id", h.handleGetEvent)

	router.Post("/", h.handleCreate)
	router.Put("/:id", h.handleUpdate)
}

func (h *EventHandler) handleCreate(ctx *fiber.Ctx) error {
	return fiber.ErrNotImplemented
}

func (h *EventHandler) handleGetManyEvents(ctx *fiber.Ctx) error {
	// limit := ctx.QueryInt("limit", 5)
	// offset := ctx.QueryInt("offset", 0)

	loc, _ := time.LoadLocation("Europe/Moscow")
	_ = time.Now().In(loc)

	from := ctx.Query("fromDate")
	to := ctx.Query("toDate")

	if from == "" || to == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(shared.Error{
			Status: fiber.StatusBadRequest,
			Code:   "INVALID_PARAMETER",
			Title:  "URL Parameter 'fromDate' or 'toDate' is missing.",
		})
	}

	fromDate, err := time.Parse(time.DateOnly, from)
	if err != nil {
		fmt.Println(err)
		return ctx.Status(fiber.StatusBadRequest).JSON(shared.Error{
			Status: fiber.StatusBadRequest,
			Code:   "INVALID_PARAMETER",
			Title:  "URL Parameter 'fromDate' is invalid.",
			Detail: "URL Parameter 'fromDate' must be date only format.",
		})
	}
	toDate, err := time.Parse(time.DateOnly, to)
	if err != nil {
		fmt.Println(err)
		return ctx.Status(fiber.StatusBadRequest).JSON(shared.Error{
			Status: fiber.StatusBadRequest,
			Code:   "INVALID_PARAMETER",
			Title:  "URL Parameter 'toDate' is invalid.",
			Detail: "URL Parameter 'toDate' must implement RFC3339 format.",
		})
	}

	if fromDate.Unix() > toDate.Unix() {
		return ctx.Status(fiber.StatusBadRequest).JSON(shared.Error{
			Status: fiber.StatusBadRequest,
			Code:   "INVALID_PARAMETER",
			Title:  "URL Parameter 'toDate' is after 'fromDate'.",
		})
	}

	log.Println(fromDate, toDate)

	events, err := h.eventRepository.FindByDate(ctx.Context(), fromDate, toDate)
	if err != nil {
		log.Println(err)
		return ctx.Status(500).JSON(shared.Error{
			Status: fiber.StatusInternalServerError,
			Code:   "INTERNAL_SERVER_ERROR",
			Title:  "Internal Server Error",
		})
	}
	return ctx.JSON(events)
}

func (h *EventHandler) handleUpdate(ctx *fiber.Ctx) error {
	return fiber.ErrNotImplemented
}

func (h *EventHandler) handleGetEvent(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id", 0)
	if err != nil {
		return err
	}

	if id == 0 {
		return shared.Error{
			Status: fiber.StatusBadRequest,
			Code:   "INVALID_PARAMETER",
			Title:  "Path parameter 'id' must be a number, greater than 0.",
		}
	}

	event, err := h.eventRepository.FindById(ctx.Context(), id)
	if err != nil {
		return err
	}

	return ctx.JSON(event)
}
