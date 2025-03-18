package events

import (
	"time"

	"github.com/gofiber/fiber/v2"

	"parteez/internal/domain/events"
	"parteez/internal/errors"
)

func NewEventHandler(eventRepository events.EventRepository, crudService events.EventCrudService) func(router fiber.Router) {
	return func(router fiber.Router) {
		router.Get("/", handleGetEvents(eventRepository))
		router.Get("/:id", handleGetEvent(eventRepository))
		router.Post("/", handleCreateEvent())
		router.Put("/:id", handleUpdateEvent())
	}
}

func handleCreateEvent() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		return fiber.ErrNotImplemented
	}
}

func handleGetEvents(events events.EventRepository) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// limit := ctx.QueryInt("limit", 5)
		// offset := ctx.QueryInt("offset", 0)

		loc, _ := time.LoadLocation("Europe/Moscow")
		_ = time.Now().In(loc)

		from := ctx.Query("fromDate")
		to := ctx.Query("toDate")

		if from == "" || to == "" {
			return errors.NewHTTPError(fiber.StatusBadRequest, errors.ErrorCodeParameterMissing, "Query parameter 'fromDate' or 'toDate' is missing.")
		}

		fromDate, err := time.Parse(time.DateOnly, from)
		if err != nil {
			return errors.NewHTTPError(fiber.StatusBadRequest, errors.ErrorCodeParameterInvalidDate, "Query parameter 'fromDate' is invalid.", "Query parameter 'fromDate' must implement RFC3339 format.")
		}
		toDate, err := time.Parse(time.DateOnly, to)
		if err != nil {
			return errors.NewHTTPError(fiber.StatusBadRequest, errors.ErrorCodeParameterInvalidDate, "Query parameter 'toDate' is invalid.", "Query parameter 'toDate' must implement RFC3339 format.")
		}

		response, err := events.FindByDate(ctx.Context(), fromDate, toDate)
		if err != nil {
			// if errors.Is(err, events.ErrEventDateRangeInvalid) {
			// 	return handler.NewHTTPError(fiber.StatusBadRequest, handler.ErrorCodeDateRangeInvalid, "Query parameter 'toDate' is after 'fromDate'.")
			// }
			return err
		}
		return ctx.JSON(response)
	}
}

func handleUpdateEvent() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		return fiber.ErrNotImplemented
	}
}

func handleGetEvent(events events.EventRepository) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		id, err := ctx.ParamsInt("id", 0)
		if err != nil {
			return err
		}

		if id == 0 {
			return errors.NewHTTPError(fiber.StatusBadRequest, errors.ErrorCodeParameterInvalidInteger, "Path parameter 'id' must be a number, greater than 0.")
		}

		response, err := events.FindById(ctx.Context(), id)
		if err != nil {
			return err
		}

		return ctx.JSON(response)
	}
}
