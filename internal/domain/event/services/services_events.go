package services

import (
	"context"
	"parteez/internal/domain/artwork"
	"parteez/internal/domain/event"
	"parteez/internal/domain/venue"
	"time"
)

type EventService struct {
	events   event.EventRepository
	genres   event.GenreRepository
	venues   venue.VenueRepository
	artworks artwork.ArtworkRepository
}

func (service *EventService) CreateDraft(ctx context.Context) (event.EventID, error) {

	now := time.Now()
	date, err := event.NewDate(now, now.Add(time.Hour*1))
	if err != nil {
		return 0, err
	}

	eventId, err := event.NewEventID(0)
	if err != nil {
		return 0, err
	}

	event, err := event.NewEvent(eventId, "title", "description", date)
	if err != nil {
		return 0, err
	}

	if err := service.events.Save(ctx, *event); err != nil {
		return 0, err
	}

	return event.ID, nil
}

func (service *EventService) Publish(id event.EventID) error {
	return nil
}
