package event

import (
	"context"
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/oklog/ulid/v2"
)

type Controller struct {
	router  fiber.Router
	storage *Storage
	log     *slog.Logger
}

func NewController(log *slog.Logger, router fiber.Router, s *Storage) *Controller {
	return &Controller{
		router:  router,
		storage: s,
		log:     log,
	}
}

func (c *Controller) Init() {
	events := c.router.Group("/events")
	events.Get("/", func(ctx *fiber.Ctx) error { return ctx.SendStatus(fiber.StatusNotImplemented) })
	events.Post("/", handleEventCreation(c.log, c.storage))
	events.Get("/search", func(ctx *fiber.Ctx) error { return ctx.SendStatus(fiber.StatusNotImplemented) })
	events.Get("/:id", func(ctx *fiber.Ctx) error { return ctx.SendStatus(fiber.StatusNotImplemented) })
	events.Patch("/:id", func(ctx *fiber.Ctx) error { return ctx.SendStatus(fiber.StatusNotImplemented) })
}

type eventCreationRequest struct {
	Name        string              `json:"name"`
	ImageURL    string              `json:"image_url,omitempty"`
	Description string              `json:"description,omitempty"`
	Genres      []string            `json:"genres,omitempty"`
	LineUp      map[string][]string `json:"line_up,omitempty"`
	StartsAt    time.Time           `json:"starts_at,omitempty"`
	EndsAt      time.Time           `json:"ends_at,omitempty"`
	MinAge      int                 `json:"min_age,omitempty"`

	TicketsURL string         `json:"tickets_url,omitempty"`
	Price      map[string]int `json:"price,omitempty"`

	LocationID int    `json:"location_id,omitempty"`
	Promoter   string `json:"promoter,omitempty"`
}

func (r eventCreationRequest) Validate() error {
	return nil
}

type eventSaver interface {
	Save(ctx context.Context, event Event) (string, error)
}

func handleEventCreation(log *slog.Logger, eventSaver eventSaver) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		const op = "event.handleCreateEvent"

		log := log.With(
			slog.String("op", op),
		)

		request := eventCreationRequest{}
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

		event := Event{
			ID:        ulid.Make(),
			Name:      request.Name,
			Status:    StatusEditing,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),

			ImageURL:    request.ImageURL,
			Description: request.Description,
			Genres:      request.Genres,
			LineUp:      lineup(request.LineUp),
			StartTime:   time.Time{},
			EndTime:     time.Time{},
			MinAge:      18,
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

func lineup(rawData map[string][]string) map[string][]Artist {
	l := make(map[string][]Artist, len(rawData))
	for i, artists := range rawData {
		stage := make([]Artist, len(artists))
		for _, name := range artists {
			stage = append(stage, Artist{Name: name})
		}
		l[i] = stage
	}
	return l
}
