package events

import (
	"context"
	"parteez/internal/domain/shared"
	"time"
)

type EventUpdate struct {
}

type EventService interface {
	CreateDraft(ctx context.Context) (*Event, error)
	Update(ctx context.Context, id EventID, update EventUpdate) (*Event, error)
	Publish(ctx context.Context, id EventID) error
	Find(ctx context.Context, from, to time.Time, page shared.Page) ([]*Event, error)
}
