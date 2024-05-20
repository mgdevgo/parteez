package storage

import (
	"context"

	"iditusi/pkg/core"
	"iditusi/pkg/storage/postgres"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Events interface {
	Save(ctx context.Context, event core.Event) (core.Event, error)
	FindByID(ctx context.Context, id int) (core.Event, error)
	Find(ctx context.Context, params core.EventSearchParams) ([]core.Event, error)
}

type Page struct {
	Limit  int
	Offset int
}

type Locations interface {
	Save(ctx context.Context, location core.Location) (core.Location, error)
	SaveAll(ctx context.Context, locations []core.Location) error
	FindByID(ctx context.Context, id int) (core.Location, error)
	FindByName(ctx context.Context, name string) (core.Location, error)
	FindAll(ctx context.Context) ([]core.Location, error)
	Delete(ctx context.Context, id int) error
}

var _ Locations = (*postgres.LocationStorage)(nil)
var _ Events = (*postgres.EventStorage)(nil)

type Storage struct {
	Events    Events
	Locations Locations
}

func NewStorage(db *pgxpool.Pool) *Storage {
	return &Storage{
		Events:    postgres.NewEventStorage(db),
		Locations: postgres.NewLocationsStorage(db),
	}
}
