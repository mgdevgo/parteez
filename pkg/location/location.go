package location

import (
	"strings"
	"time"
)

type Kind string

const (
	Unknown     Kind = "UNKNOWN"
	Club        Kind = "CLUB"
	Bar         Kind = "BAR"
	Cafe        Kind = "CAFE"
	ConcertHall Kind = "CONCERT_HALL"
	Space       Kind = "SPACE"
)

func parseLocationType(value string) Kind {
	switch strings.ToLower(value) {
	case "club":
		return Club
	case "bar":
		return Bar
	case "cafe":
		return Cafe
	case "concert_hall":
		return ConcertHall
	case "space":
		return Space
	default:
		return Unknown
	}
}

type Location struct {
	ID          int
	Name        string
	Type        Kind
	Description string
	ArtworkURL  string
	// IsMultiStage bool
	Stages        []string
	Address       string
	MetroStations []string
	// TODO
	// Latitude  string
	// Longitude string
	IsPublic  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewLocation(name string) Location {
	// if len(stages) <= 0 {
	// 	stages = StagesDefault
	// } else {
	// 	for i := 0; i < len(stages); i++ {
	// 		stages[i] = strings.ToUpper(stages[i])
	// 	}
	// }
	return Location{
		Name:   name,
		Type:   Unknown,
		Stages: []string{"main"},
	}
}

func (l *Location) IsEmpty() bool {
	if l.ID == 0 {
		return true
	}

	// Todo: check other values
	return false
}

var StagesDefault = []string{"main"}
