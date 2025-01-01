package storage

import (
	"context"
	"time"

	"iditusi/internal/models"
)

type BaseStorage[T any] interface {
	Save(ctx context.Context, entity T) (T, error)
	FindAll(ctx context.Context) ([]T, error)
	FindById(ctx context.Context, id int) (T, error)
	Delete(ctx context.Context, id int) error
}

type EventStorage interface {
	BaseStorage[models.Event]
	FindByDate(ctx context.Context, fromDate time.Time, toDate time.Time) ([]models.Event, error)
	// Publish(ctx context.Context, eventId int) error
}

type VenueStorage interface {
	BaseStorage[models.Venue]
	FindByName(ctx context.Context, name string) (models.Venue, error)
	// Show(ctx context.Context, locationId int) error
}

type UserStorage interface {
	BaseStorage[models.User]
	SetRole(ctx context.Context, userId int, roleName string) error
}

type LineUpStorage interface {
	BaseStorage[models.LineUp]
}

