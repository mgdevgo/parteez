package repositories

import (
	"context"
	"time"

	"iditusi/internal/models"
)

type BaseRepository[T any] interface {
	Save(ctx context.Context, entity T) (T, error)
	FindAll(ctx context.Context) ([]T, error)
	FindById(ctx context.Context, id int) (T, error)
	Delete(ctx context.Context, id int) error
}

type EventRepository interface {
	BaseRepository[models.Event]
	FindByDate(ctx context.Context, fromDate time.Time, toDate time.Time) ([]models.Event, error)
	// Publish(ctx context.Context, eventId int) error
}

type VenueRepository interface {
	BaseRepository[models.Venue]
	FindByName(ctx context.Context, name string) (models.Venue, error)
	// Show(ctx context.Context, locationId int) error
}

type UserRepository interface {
	BaseRepository[models.User]
	SetRole(ctx context.Context, userId int, roleName string) error
}
