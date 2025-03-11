package venue

import (
	"context"

	"parteez/internal/repository"
)

type VenueRepository interface {
	repository.Repository[*Venue]
	FindByName(ctx context.Context, name string) (*Venue, error)
}
