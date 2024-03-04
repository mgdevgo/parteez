package location

import (
	"strings"
	"time"
)

type LocationType string

const (
	LocationUnknown     LocationType = "UNKNOWN"
	LocationClub        LocationType = "CLUB"
	LocationBar         LocationType = "BAR"
	LocationCafe        LocationType = "CAFE"
	LocationConcertHall LocationType = "CONCERT_HALL"
	LocationSpace       LocationType = "SPACE"
)

func parseLocationType(value string) LocationType {
	switch strings.ToLower(value) {
	case "club":
		return LocationClub
	case "bar":
		return LocationBar
	case "cafe":
		return LocationCafe
	case "concert_hall":
		return LocationConcertHall
	case "space":
		return LocationSpace
	default:
		return LocationUnknown
	}
}

type Location struct {
	ID            string
	Name          string
	Type          LocationType
	Description   string
	ImageURL      string
	MusicGenres   []string
	Stages        []string
	Address       string
	MetroStations []string
	// TODO
	// Latitude  string
	// Longitude string
	Public    bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewLocation(id string, name string, locationType LocationType, description string, address string, stages []string) Location {
	if len(stages) <= 0 {
		stages = StagesDefault
	} else {
		for i := 0; i < len(stages); i++ {
			stages[i] = strings.ToUpper(stages[i])
		}
	}
	return Location{
		ID:          id,
		Name:        name,
		Type:        locationType,
		Description: description,
		Stages:      stages,
		Address:     address,
		Public:      false,
	}
}

var StagesDefault = []string{"main"}
