package service

import (
	"context"
	"time"

	"parteez/internal/domain/artwork"
	"parteez/internal/domain/events"
	"parteez/internal/domain/venue"
	"parteez/internal/page"
)

type EventCrudService struct {
	events   events.EventRepository
	venues   venue.VenueRepository
	artworks artwork.ArtworkRepository
}

func NewEventCrudService(events events.EventRepository, venues venue.VenueRepository, artworks artwork.ArtworkRepository) *EventCrudService {
	return &EventCrudService{
		events:   events,
		venues:   venues,
		artworks: artworks,
	}
}

func (service *EventCrudService) CreateDraft(ctx context.Context) (*events.Event, error) {

	now := time.Now()
	date, err := events.NewDate(now, now.Add(time.Hour*1))
	if err != nil {
		return nil, err
	}

	eventId, err := events.NewEventID(0)
	if err != nil {
		return nil, err
	}

	event, err := events.NewEvent(eventId, "title", "description", date)
	if err != nil {
		return nil, err
	}

	if err := service.events.Save(ctx, event); err != nil {
		return nil, err
	}

	return event, nil
}

func (service *EventCrudService) Publish(ctx context.Context, id events.EventID) error {
	panic("not implemented")
}

func (service *EventCrudService) Find(ctx context.Context, from, to time.Time, page page.Page) ([]*events.Event, error) {
	panic("not implemented")
}

func (service *EventCrudService) Update(ctx context.Context, id events.EventID, update events.EventUpdate) (*events.Event, error) {
	panic("not implemented")
}
