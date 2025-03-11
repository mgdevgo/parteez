package events

import (
	"context"
	"time"

	"parteez/internal/repository"
)

type EventRepository interface {
	repository.Repository[*Event]
	FindByDate(ctx context.Context, fromDate time.Time, toDate time.Time) ([]*Event, error)
	// Publish(ctx context.Context, eventId int) error
}
