package event

import (
	"context"
	"time"

	"parteez/internal/domain/shared"
)

type EventRepository interface {
	shared.Repository[Event]
	FindByDate(ctx context.Context, fromDate time.Time, toDate time.Time) ([]Event, error)
	// Publish(ctx context.Context, eventId int) error
}

type GenreRepository interface {
	shared.Repository[Genre]
}
