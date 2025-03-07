package server

import (
	event "parteez/internal/domain/events"
	"parteez/internal/domain/shared"
	"parteez/internal/domain/venue"
	venueHandler "parteez/internal/domain/venue/handler"
)

type Handlers struct {
	health *shared.HealthHandler
	events *event.EventHandler
	venues *venueHandler.VenueHandler
}

func NewHandlers(eventRepository event.EventRepository, eventCrudService event.EventCrudService, venueRepository venue.VenueRepository) *Handlers {
	return &Handlers{
		health: shared.NewHealthHandler(),
		events: event.NewEventHandler(eventRepository, eventCrudService),
		venues: venueHandler.NewVenueHandler(venueRepository),
	}
}
