package listing

import (
	"context"

	"iditusi/pkg/core"
)

type EventStorage interface {
	FindEventById(ctx context.Context, id int) (core.Event, error)
	FindEvents(ctx context.Context, fromDate string, toDate string) ([]core.Event, error)
}

type EventListingLocalService struct {
	eventStorage EventStorage
}

func NewEventListingService(eventStorage EventStorage) *EventListingLocalService {
	return &EventListingLocalService{
		eventStorage: eventStorage,
	}
}
