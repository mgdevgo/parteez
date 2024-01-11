package event

import (
	"context"
	"log/slog"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/oklog/ulid/v2"
)

type EventRouterConfig struct{}

func EventRouter(log *slog.Logger, router fiber.Router, storage *Storage) {
	events := router.Group("/events")
	events.Get("/", handleGetEvent(log, storage))
	events.Post("/", handleCreateEvent(log, storage))
	events.Get("/search", func(ctx *fiber.Ctx) error { return ctx.SendStatus(fiber.StatusNotImplemented) })
	events.Get("/:id", func(ctx *fiber.Ctx) error { return ctx.SendStatus(fiber.StatusNotImplemented) })
	events.Patch("/:id", func(ctx *fiber.Ctx) error { return ctx.SendStatus(fiber.StatusNotImplemented) })
}

type EventGetter interface {
	Get(ctx context.Context, eventIDs []string) ([]Event, error)
}

type EventResponse struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	UpdateTime time.Time `json:"update_time"`
}

func handleGetEvent(log *slog.Logger, eventgetter EventGetter) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		ids := ctx.Query("ids")
		idsArray := strings.Split(strings.Trim(ids, "[]"), ",")

		events, err := eventgetter.Get(ctx.Context(), idsArray)
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		response := make([]EventResponse, 0, len(events))
		for _, e := range events {
			item := EventResponse{
				ID:         e.ID.String(),
				Name:       e.Name,
				UpdateTime: e.UpdatedAt,
			}
			response = append(response, item)
		}

		return ctx.JSON(response)
	}
}

type EventCreationRequest struct {
	Name        string            `json:"name"`
	ImageURL    string            `json:"image_url,omitempty"`
	Description string            `json:"description,omitempty"`
	Genres      []string          `json:"genres,omitempty"`
	LineUp      map[string]LineUp `json:"line_up,omitempty"`
	StartTime   time.Time         `json:"start_time,omitempty"`
	EndTime     time.Time         `json:"end_time,omitempty"`
	MinAge      int               `json:"min_age,omitempty"`

	TicketsURL string         `json:"tickets_url,omitempty"`
	Price      map[string]int `json:"price,omitempty"`

	LocationID int    `json:"location_id,omitempty"`
	Promoter   string `json:"promoter,omitempty"`
}

func (r EventCreationRequest) Validate() error {
	return nil
}

type eventSaver interface {
	Save(ctx context.Context, event Event) (string, error)
}

func handleCreateEvent(log *slog.Logger, eventSaver eventSaver) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		const op = "event.handleCreateEvent"

		log := log.With(
			slog.String("op", op),
		)

		var request EventCreationRequest
		if err := c.BodyParser(&request); err != nil {
			log.Error("failed to process request data", slog.Attr{
				Key:   "error",
				Value: slog.StringValue(err.Error()),
			})
			// return api.ErrorResponse(api.NewError())
			return c.SendStatus(500)
		}

		if err := request.Validate(); err != nil {
			return c.SendStatus(400)
		}

		// TODO: locationID.isValid()
		event := Event{
			ID:        ulid.Make(),
			Name:      request.Name,
			Status:    StatusEditing,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),

			ImageURL:    request.ImageURL,
			Description: request.Description,
			Genres:      request.Genres,
			LineUp:      request.LineUp,
			StartTime:   request.StartTime,
			EndTime:     request.EndTime,
			MinAge:      DefaultMinAge,
			TicketsURL:  request.TicketsURL,
			Price:       request.Price,
			LocationID:  request.LocationID,
			Promoter:    request.Promoter,
		}

		if _, err := eventSaver.Save(context.TODO(), event); err != nil {
			log.Error("failed to process use case", slog.Attr{
				Key:   "error",
				Value: slog.StringValue(err.Error()),
			})
			return c.SendStatus(500)
		}

		return c.SendStatus(201)
	}
}

//func lineup(rawData map[string][]string) map[string][]Artist {
//	l := make(map[string][]Artist)
//	for i, artists := range rawData {
//		stage := make([]Artist, 0, len(artists))
//		for _, name := range artists {
//			stage = append(stage, Artist{Name: name})
//		}
//		l[i] = stage
//	}
//	return l
//}
