package event

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// import (
//
//	"context"
//	"log/slog"
//	"strings"
//	"time"
//
//	"github.com/gofiber/fiber/v2"
//	"github.com/oklog/ulid/v2"
//
// )
//
// type EventRouterConfig struct{}
//
//	func EventRouter(log *slog.Logger, router fiber.Router, storage EventStorage) {
//		events := router.Group("/events")
//		events.Get("/", handleGetEvent(log, storage))
//		events.Post("/", handleCreateEvent(log, storage))
//		events.Get("/search", func(ctx *fiber.Ctx) error { return ctx.SendStatus(fiber.StatusNotImplemented) })
//		events.Get("/:id", func(ctx *fiber.Ctx) error { return ctx.SendStatus(fiber.StatusNotImplemented) })
//		events.Patch("/:id", func(ctx *fiber.Ctx) error { return ctx.SendStatus(fiber.StatusNotImplemented) })
//	}
//
//	type EventGetter interface {
//		Get(ctx context.Context, eventIDs []string) ([]Event, error)
//	}

type handler struct {
	eventService *service
}

func NewHandler(eventService *service) *handler {
	return &handler{
		eventService: eventService,
	}
}

func (h *handler) RegisterRoutes(router fiber.Router) {
	events := router.Group("/events")
	events.Get("/", h.ListEvents)
	events.Post("/", func(ctx *fiber.Ctx) error { return ctx.SendStatus(fiber.StatusNotImplemented) })
	events.Get("/search", func(ctx *fiber.Ctx) error { return ctx.SendStatus(fiber.StatusNotImplemented) })
	events.Get("/:id", func(ctx *fiber.Ctx) error { return ctx.SendStatus(fiber.StatusNotImplemented) })
	events.Patch("/:id", func(ctx *fiber.Ctx) error { return ctx.SendStatus(fiber.StatusNotImplemented) })

}

func (h *handler) ListEvents(ctx *fiber.Ctx) error {
	limit := ctx.Query("limit")
	offset := ctx.Query("offset")

	var request ListEventsRequest
	var err error

	request.Limit, err = strconv.Atoi(limit)
	if err != nil {
		return err
	}
	request.Offset, err = strconv.Atoi(offset)
	if err != nil {
		return err
	}

	response, err := h.eventService.ListEvents(request)
	// if err != nil {
	// 	return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
	// 		"error": err.Error(),
	// 	})
	// }
	//
	// response := make([]EventResponse, 0, len(events))
	// for _, e := range events {
	// 	item := EventResponse{
	// 		ID:         e.ID.String(),
	// 		Name:       e.Name,
	// 		UpdateTime: e.UpdatedAt,
	// 	}
	// 	response = append(response, item)
	// }

	// originalURL := ctx.OriginalURL()
	// path := strings.Split(originalURL, "?")[0]
	// next := fmt.Sprintf("%s?offset=%s", url.PathEscape(path), url.PathEscape(offset))

	return ctx.JSON(response)
}

//
// type EventCreationRequest struct {
// 	Name        string                  `json:"name"`
// 	LogoURL    string                  `json:"image_url,omitempty"`
// 	Description string                  `json:"description,omitempty"`
// 	Genres      []string                `json:"genres,omitempty"`
// 	LineUp      map[string]event.LineUp `json:"line_up,omitempty"`
// 	StartDate   time.Time               `json:"start_time,omitempty"`
// 	EndDate     time.Time               `json:"end_time,omitempty"`
// 	MinAge      int                     `json:"min_age,omitempty"`
//
// 	TicketsURL string         `json:"tickets_url,omitempty"`
// 	Price      map[string]int `json:"price,omitempty"`
//
// 	LocationID int    `json:"location_id,omitempty"`
// 	Promoter   string `json:"promoter,omitempty"`
// }
//
// func (r EventCreationRequest) Validate() error {
// 	return nil
// }
//
// type eventSaver interface {
// 	Save(ctx context.Context, event event.Event) (string, error)
// }
//
// func handleCreateEvent(log *slog.Logger, eventSaver eventSaver) func(c *fiber.Ctx) error {
// 	return func(c *fiber.Ctx) error {
// 		const op = "event.handleCreateEvent"
//
// 		log := log.With(
// 			slog.String("op", op),
// 		)
//
// 		var request EventCreationRequest
// 		if err := c.BodyParser(&request); err != nil {
// 			log.Error("failed to process request data", slog.Attr{
// 				Key:   "error",
// 				Value: slog.StringValue(err.Error()),
// 			})
// 			// return api.ErrorResponse(api.NewError())
// 			return c.SendStatus(500)
// 		}
//
// 		if err := request.Validate(); err != nil {
// 			return c.SendStatus(400)
// 		}
//
// 		// TODO: locationID.isValid()
// 		event := event.Event{
// 			ID:        ulid.Make(),
// 			Name:      request.Name,
// 			Status:    event.StatusEditing,
// 			CreatedAt: time.Now(),
// 			UpdatedAt: time.Now(),
//
// 			LogoURL:    request.LogoURL,
// 			Description: request.Description,
// 			Genres:      request.Genres,
// 			LineUp:      request.LineUp,
// 			StartDate:   request.StartDate,
// 			EndDate:     request.EndDate,
// 			MinAge:      event.DefaultMinAge,
// 			TicketsURL:  request.TicketsURL,
// 			Price:       request.Price,
// 			LocationID:  request.LocationID,
// 			Promoter:    request.Promoter,
// 		}
//
// 		if _, err := eventSaver.Save(context.TODO(), event); err != nil {
// 			log.Error("failed to process use case", slog.Attr{
// 				Key:   "error",
// 				Value: slog.StringValue(err.Error()),
// 			})
// 			return c.SendStatus(500)
// 		}
//
// 		return c.SendStatus(201)
// 	}
// }

// func lineup(rawData map[string][]string) map[string][]Artist {
//	l := make(map[string][]Artist)
//	for i, artists := range rawData {
//		stage := make([]Artist, 0, len(artists))
//		for _, name := range artists {
//			stage = append(stage, Artist{Name: name})
//		}
//		l[i] = stage
//	}
//	return l
// }
