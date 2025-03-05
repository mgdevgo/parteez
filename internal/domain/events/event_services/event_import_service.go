package event_services

import (
	"context"

	"parteez/internal/domain/events"
	"parteez/internal/domain/events/scraping"
)

type EventImportService struct {
	scrapingService *scraping.ScrapingService
	crudService     *events.EventCrudService
}

func NewEventImportService(scrapingService *scraping.ScrapingService, crudService *events.EventCrudService) *EventImportService {
	return &EventImportService{
		scrapingService: scrapingService,
		crudService:     crudService,
	}
}

func (service *EventImportService) Import(ctx context.Context) error {
	panic("not implemented")
}
