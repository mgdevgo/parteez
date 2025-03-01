package venue

import (
	"time"

	"parteez/internal/domain/artwork"
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

type VenueVisability string

const (
	VenueVisabilityPublic  VenueVisability = "PUBLIC"
	VenueVisabilityPrivate VenueVisability = "PRIVATE"
)

type VenueID int

type Venue struct {
	ID          VenueID
	Name        string
	Description string
	Type        VenueType
	Artwork     artwork.ArtworkID
	Stages      []string
	Location    Location
	Visability  VenueVisability
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewVeue(id VenueID, name string, description string) Venue {
	return Venue{
		ID:          0,
		Name:        "",
		Description: "",
		Type:        VenueTypeDefault,
		Artwork:     0,
		Stages:      []string{VENUE_DEFAULT_STAGE},
		Location: Location{
			Address:       "",
			MetroStations: make([]string, 0),
			Latitude:      "",
			Longitude:     "",
		},
		Visability: VenueVisabilityPublic,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
}

type Location struct {
	Address       string
	MetroStations []string
	Latitude      string
	Longitude     string
}
