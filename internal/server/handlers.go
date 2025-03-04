package server

import (
	event "parteez/internal/domain/events/handlers"
	"parteez/internal/domain/shared"
	venue "parteez/internal/domain/venue/handlers"
)

type Handlers struct {
	health *shared.HealthHandler
	events *event.EventHandler
	venues *venue.VenueHandler
}

func NewHandlers() *Handlers {
	return &Handlers{
		health: shared.NewHealthHandler(),
		events: event.NewEventHandler(),
		venues: venue.NewVenueHandler(),
	}
}
