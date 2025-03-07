package venue

import (
	"errors"
	"fmt"
	"time"

	"parteez/internal/domain/artwork"
)

var (
	ErrVenue                  = errors.New("venue")
	ErrVenueNotFound          = fmt.Errorf("%w: not found", ErrVenue)
	ErrVenueAlreadyExists     = fmt.Errorf("%w: already exists", ErrVenue)
	ErrVenueInvalidID         = fmt.Errorf("%w: invalid id", ErrVenue)
	ErrVenueInvalidVisability = fmt.Errorf("%w: invalid visability state", ErrVenue)
	ErrVenueInvalidType       = fmt.Errorf("%w: invalid type", ErrVenue)
)

const VENUE_DEFAULT_STAGE = "main"

type VenueType string

const (
	VenueTypeDefault     VenueType = "DEFAULT"
	VenueTypeClub        VenueType = "CLUB"
	VenueTypeBar         VenueType = "BAR"
	VenueTypeCafe        VenueType = "CAFE"
	VenueTypeConcertHall VenueType = "CONCERT_HALL"
	VenueTypeSpace       VenueType = "SPACE"
)

var venueTypes = map[VenueType]struct{}{
	VenueTypeDefault:     {},
	VenueTypeClub:        {},
	VenueTypeBar:         {},
	VenueTypeCafe:        {},
	VenueTypeConcertHall: {},
	VenueTypeSpace:       {},
}

func NewVenueType(venueType string) (VenueType, error) {
	if _, ok := venueTypes[VenueType(venueType)]; !ok {
		return "", fmt.Errorf("%w: venueType=%s", ErrVenueInvalidType, venueType)
	}
	return VenueType(venueType), nil
}

type VenueVisability string

const (
	VenueVisabilityPublic  VenueVisability = "PUBLIC"
	VenueVisabilityPrivate VenueVisability = "PRIVATE"
)

func NewVenueVisability(visability string) (VenueVisability, error) {
	switch visability {
	case "PUBLIC":
		return VenueVisabilityPublic, nil
	case "PRIVATE":
		return VenueVisabilityPrivate, nil
	default:
		return "", fmt.Errorf("%w: visability=%s", ErrVenueInvalidVisability, visability)
	}
}

type VenueID int

func NewVenueID(id int) (VenueID, error) {
	if id < 0 {
		return 0, fmt.Errorf("%w: id=%d", ErrVenueInvalidID, id)
	}
	return VenueID(id), nil
}

type Venue struct {
	ID          VenueID
	Name        string
	Description string
	Type        VenueType
	ArtworkID   artwork.ArtworkID
	Stages      []string
	Location    Location
	Visability  VenueVisability
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewVeue(id VenueID, name string, description string, venueType VenueType, stages []string, location Location, visability VenueVisability, createdAt time.Time, updatedAt time.Time) *Venue {
	return &Venue{
		ID:          id,
		Name:        name,
		Description: description,
		Type:        venueType,
		ArtworkID:   0,
		Stages:      []string{VENUE_DEFAULT_STAGE},
		Location:    location,
		Visability:  visability,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}
}

type Location struct {
	Address       string
	MetroStations []string
	Latitude      string
	Longitude     string
}

func NewLocation(address string, metroStations []string, latitude string, longitude string) Location {
	return Location{
		Address:       address,
		MetroStations: metroStations,
		Latitude:      latitude,
		Longitude:     longitude,
	}
}
