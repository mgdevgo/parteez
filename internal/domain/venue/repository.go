package venue

import (
	"context"
	"parteez/internal/domain/shared"
)

type VenueRepository interface {
	shared.Repository[*Venue]
	FindByName(ctx context.Context, name string) (*Venue, error)
	// Show(ctx context.Context, locationId int) error
}
