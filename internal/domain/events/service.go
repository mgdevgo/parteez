package events

import (
	"context"
	"time"

	"parteez/internal/domain/shared"
)

type EventUpdate struct {
}

type EventCrudService interface {
	CreateDraft(ctx context.Context) (*Event, error)
	Update(ctx context.Context, id EventID, update EventUpdate) (*Event, error)
	Publish(ctx context.Context, id EventID) error
	Find(ctx context.Context, from, to time.Time, page shared.Page) ([]*Event, error)
}
