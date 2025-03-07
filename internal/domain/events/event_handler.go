package events

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"

	"parteez/internal/domain/shared/handler"
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
	router.Get("/", h.handleGetEvents)
	router.Get("/:id", h.handleGetEvent)

	router.Post("/", h.handleCreate)
	router.Put("/:id", h.handleUpdate)
}

func (h *EventHandler) handleCreate(ctx *fiber.Ctx) error {
	return fiber.ErrNotImplemented
}

func (h *EventHandler) handleGetEvents(ctx *fiber.Ctx) error {
	// limit := ctx.QueryInt("limit", 5)
	// offset := ctx.QueryInt("offset", 0)

	loc, _ := time.LoadLocation("Europe/Moscow")
	_ = time.Now().In(loc)

	from := ctx.Query("fromDate", "")
	to := ctx.Query("toDate", "")

	if from == "" || to == "" {
		return handler.NewHTTPError(fiber.StatusBadRequest, handler.ErrorCodeParameterMissing, "Query parameter 'fromDate' or 'toDate' is missing.")
	}

	fromDate, err := time.Parse(time.DateOnly, from)
	if err != nil {
		fmt.Println(err)
		return handler.NewHTTPError(fiber.StatusBadRequest, handler.ErrorCodeParameterInvalidStringBlank, "Query parameter 'fromDate' is invalid.", "URL Parameter 'fromDate' must implement RFC3339 format.")
	}
	toDate, err := time.Parse(time.DateOnly, to)
	if err != nil {
		fmt.Println(err)
		return handler.NewHTTPError(fiber.StatusBadRequest, handler.ErrorCodeParameterInvalidStringBlank, "Query parameter 'toDate' is invalid.", "URL Parameter 'toDate' must implement RFC3339 format.")
	}

	events, err := h.eventRepository.FindByDate(ctx.Context(), fromDate, toDate)
	if err != nil {
		// if errors.Is(err, events.ErrEventDateRangeInvalid) {
		// 	return handler.NewHTTPError(fiber.StatusBadRequest, handler.ErrorCodeDateRangeInvalid, "Query parameter 'toDate' is after 'fromDate'.")
		// }
		return err
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
		return handler.NewHTTPError(fiber.StatusBadRequest, handler.ErrorCodeParameterInvalidInteger, "Path parameter 'id' must be a number, greater than 0.")
	}

	event, err := h.eventRepository.FindById(ctx.Context(), id)
	if err != nil {
		return err
	}

	return ctx.JSON(event)
}
