package events

import (
	"context"
	"time"

	"iditusi/pkg/core"
	"iditusi/pkg/core/repository"
)

type EventRepository interface {
	repository.CrudRepository[core.Event]
	FindByDate(ctx context.Context, fromDate time.Time, toDate time.Time) ([]core.Event, error)
}

// type EventRepository interface {
// 	Save(ctx context.Context, event core.Event) (core.Event, error)
// 	Find(ctx context.Context, params core.EventSearchParams) ([]core.Event, error)
// 	FindById(ctx context.Context, id int) (core.Event, error)
// 	DeleteById(ctx context.Context, id int) error
// }
