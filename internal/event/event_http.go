package event

import (
	"context"
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v2"
)

type eventCreationRequest struct {
	Name        string              `json:"name"`
	ArtworkURL  string              `json:"artwork_url"`
	Description string              `json:"description"`
	Genres      []string            `json:"genres"`
	LineUp      map[string][]string `json:"line_up"`
	StartsAt    time.Time           `json:"starts_at"`
	EndsAt      time.Time           `json:"ends_at"`
	MinAge      int                 `json:"min_age"`

	TicketsURL string         `json:"tickets_url"`
	Price      map[string]int `json:"price"`

	Location string `json:"location"`
	Promoter string `json:"promoter"`
}

func (r eventCreationRequest) ToEvent() *Event {
	lineup := map[string][]Artist{}
	for i, artists := range r.LineUp {
		stage := make([]Artist, len(artists))
		for _, name := range artists {
			stage = append(stage, Artist{Name: name})
		}
		lineup[i] = stage
	}
	return &Event{
		Name:        r.Name,
		ArtworkURL:  r.ArtworkURL,
		Description: r.Description,
		Genres:      r.Genres,
		LineUp:      lineup,
		StartsAt:    r.StartsAt,
		EndsAt:      r.EndsAt,
		TicketsURL:  r.TicketsURL,
		Price:       r.Price,
		Location:    Location{Name: r.Name},
		Promoter:    r.Promoter,
		Status:      StatusEditing,
	}
}

func (r eventCreationRequest) Validate() error {
	return nil
}

type eventSaver interface {
	Save(ctx context.Context, event *Event) (int, error)
}

func HandleCreateEvent(log *slog.Logger, eventSaver eventSaver) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		const op = "event.HandleCreateEvent"

		log := log.With(
			slog.String("op", op),
		)

		var request eventCreationRequest
		if err := c.BodyParser(request); err != nil {
			log.Error("failed to process request data", slog.Attr{
				Key:   "error",
				Value: slog.StringValue(err.Error()),
			})
			return c.SendStatus(500)
		}

		// TODO request validation

		if _, err := eventSaver.Save(context.TODO(), request.ToEvent()); err != nil {
			log.Error("failed to process use case", slog.Attr{
				Key:   "error",
				Value: slog.StringValue(err.Error()),
			})
			return c.SendStatus(500)
		}

		return c.SendStatus(201)
	}
}

type EventResponse struct {
}
