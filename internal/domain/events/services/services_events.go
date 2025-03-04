package services

import (
	"context"
	"time"

	"parteez/internal/domain/artwork"
	"parteez/internal/domain/events"
	"parteez/internal/domain/venue"
)

type EventService struct {
	events   events.EventRepository
	genres   events.GenreRepository
	venues   venue.VenueRepository
	artworks artwork.ArtworkRepository
}

func (service *EventService) CreateDraft(ctx context.Context) (events.EventID, error) {

	now := time.Now()
	date, err := events.NewDate(now, now.Add(time.Hour*1))
	if err != nil {
		return 0, err
	}

	eventId, err := events.NewEventID(0)
	if err != nil {
		return 0, err
	}

	event, err := events.NewEvent(eventId, "title", "description", date)
	if err != nil {
		return 0, err
	}

	if err := service.events.Save(ctx, *event); err != nil {
		return 0, err
	}

	return event.ID, nil
}

func (service *EventService) Publish(id events.EventID) error {
	return nil
}
