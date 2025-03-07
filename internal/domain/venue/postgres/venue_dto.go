package postgres

import (
	"time"

	"parteez/internal/domain/venue"
)

type venueRow struct {
	ID            int
	Name          string
	Type          string
	Description   string
	ArtworkID     int
	Stages        []string
	Address       string
	MetroStations []string
	Latitude      string
	Longitude     string
	Visability    string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func venueToRow(venue *venue.Venue) venueRow {
	return venueRow{
		ID:            int(venue.ID),
		Name:          venue.Name,
		Type:          string(venue.Type),
		Description:   venue.Description,
		ArtworkID:     int(venue.ArtworkID),
		Stages:        venue.Stages,
		Address:       venue.Location.Address,
		MetroStations: venue.Location.MetroStations,
		Latitude:      venue.Location.Latitude,
		Longitude:     venue.Location.Longitude,
		Visability:    string(venue.Visability),
		CreatedAt:     venue.CreatedAt,
		UpdatedAt:     venue.UpdatedAt,
	}
}

func rowToVenue(row venueRow) (*venue.Venue, error) {
	location := venue.NewLocation(
		row.Address,
		row.MetroStations,
		row.Latitude,
		row.Longitude,
	)

	id, err := venue.NewVenueID(row.ID)
	if err != nil {
		return nil, err
	}

	visability, err := venue.NewVenueVisability(row.Visability)
	if err != nil {
		return nil, err
	}

	venueType, err := venue.NewVenueType(row.Type)
	if err != nil {
		return nil, err
	}

	return venue.NewVeue(
		id,
		row.Name,
		row.Description,
		venueType,
		row.Stages,
		location,
		visability,
		row.CreatedAt,
		row.UpdatedAt,
	), nil
}
