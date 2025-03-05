package server

import (
	event "parteez/internal/domain/events"
	"parteez/internal/domain/shared"
	venue "parteez/internal/domain/venue/handlers"
)

type Handlers struct {
	health *shared.HealthHandler
	events *event.EventHandler
	venues *venue.VenueHandler
}

func NewHandlers(eventRepository event.EventRepository, eventCrudService event.EventCrudService) *Handlers {
	return &Handlers{
		health: shared.NewHealthHandler(),
		events: event.NewEventHandler(eventRepository, eventCrudService),
		venues: venue.NewVenueHandler(),
	}
}
