package storage

import (
	"iditusi/internal/location"
	"iditusi/internal/shared/entity"

	"github.com/google/uuid"
	"github.com/oklog/ulid/v2"
)

type Entity interface {
	entity.Event | location.Location
}

type ID interface {
	int | string | uuid.UUID | ulid.ULID
}

type CRUDWrapper[TEntity Entity, TID ID] interface {
	Create(entity TEntity) (TID, error)
	Update(id TID, options map[string]any) (TEntity, error)
	Delete(id TID) error
	Transaction() (int, error)
}
